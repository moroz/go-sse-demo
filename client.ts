#!/usr/bin/env -S deno run --allow-net

const es = new EventSource("http://localhost:4000/events");
es.addEventListener("message", (msg) => {
  console.log(msg.data);
});
