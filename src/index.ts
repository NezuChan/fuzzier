import { Elysia } from "elysia";
import { createClient } from "redis";
import { Entity, EntityId, Repository, Schema } from "redis-om";
import { logger, serializers, serializeRequest } from "@bogeychan/elysia-logger";
import pretty from "pino-pretty";
import { randomUUID } from "node:crypto";

// As NezukoChan still use v3 Lavalink with branch custom, not the stable release
const TrackSchema = new Schema("track", {
    track: { type: "text" },
    info_identifier: { type: "text", path: "$.info.identifier" },
    info_isSeekable: { type: "boolean", path: "$.info.isSeekable" },
    info_author: { type: "text", path: "$.info.author" },
    info_artworkUrl: { type: "text", path: "$.info.artworkUrl" },
    info_length: { type: "number", path: "$.info.length" },
    info_isStream: { type: "boolean", path: "$.info.isStream" },
    info_position: { type: "number", path: "$.info.position" },
    info_sourceName: { type: "text", path: "$.info.sourceName" },
    info_title: { type: "text", path: "$.info.title" },
    info_uri: { type: "text", path: "$.info.uri" }
}, {
    dataStructure: "JSON"
});

const redis = createClient({
    url: process.env.REDIS_URL
});
const repository = new Repository(TrackSchema, redis);
redis.on("error", err => console.log("Redis Client Error", err));

const stream = pretty({
    colorize: true,
    translateTime: "SYS:yyyy-mm-dd HH:MM:ss.l o"
});

const app = new Elysia()
    .use(logger({
        name: "Nezu Fuzzier",
        timestamp: true,
        serializers: {
            ...serializers,
            request: (request: Request) => {
                const url = new URL(request.url);

                return {
                    ...serializeRequest(request),
                    id: request.headers.get("X-Request-ID") ?? randomUUID(),
                    path: url.pathname
                };
            }
        },
        stream
    }))
    .derive(ctx => {
        ctx.log.info(ctx.request, "Request");
        return { };
    })
    .get("/search", async ctx => {
        if (ctx.headers.authorization !== process.env.AUTHORIZATION) {
            ctx.set.status = 401;
            return { message: "Unauthorized" };
        }

        try {
            const tracks = await repository.search()
                .where("info_title")
                .match(`*${decodeURIComponent(ctx.query.q as string)}*`)
                .or("info_author")
                .match(`*${decodeURIComponent(ctx.query.q as string)}*`)
                .or("info_uri")
                .match(`*${decodeURIComponent(ctx.query.q as string)}*`)
                .return.page(0, 25);

            return tracks;
        } catch (e) {
            ctx.log.error(e, "Error");
            return [];
        }
    })
    .post("/", async ctx => {
        if (ctx.headers.authorization !== process.env.AUTHORIZATION) {
            ctx.set.status = 401;
            return { message: "Unauthorized" };
        }

        try {
            const { tracks } = ctx.body as { tracks: Entity[] };
            for (const track of tracks) {
                const entity = await repository.save(track);
                await repository.expire(entity[EntityId]!, 3 * 60 * 60);
            }

            return { message: "OK" };
        } catch (e) {
            ctx.log.error(e, "Error");
            return { message: "Error", error: (e as { message: string }).message };
        }
    })
    .listen(Number(process.env.PORT ?? 3000));


await redis.connect();
await repository.createIndex();

console.log(`🦊 Nezu Fuzzier is running at on port ${app.server!.port}...`);
