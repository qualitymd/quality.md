#!/usr/bin/env node
import { Effect } from "effect";
import { NodeRuntime, NodeServices } from "@effect/platform-node";
import { run } from "./cli.ts";

run.pipe(Effect.provide(NodeServices.layer), NodeRuntime.runMain);
