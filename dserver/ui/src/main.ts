import { RenderBuilder } from "./renderBuilder";

var renderer = new RenderBuilder(".body", 600, 16);
var ws = new WebSocket("ws://localhoast:8282");

ws.addEventListener("open", (event) => {
  console.log(event);
});

ws.addEventListener("message", (event) => {
  console.log(event);
});
