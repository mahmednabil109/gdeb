import { RenderBuilder } from "./renderBuilder";

let body = document.querySelector(".body") || { innerHTML: "" };

var renderer = new RenderBuilder(".body", 600, 12);
var ws = new WebSocket("ws://localhost:8282/");

var nodes = new Map<string, [string, string]>();

ws.addEventListener("open", (event) => {
  console.log(event);
});

ws.addEventListener("message", (event) => {
  let data = JSON.parse(event.data);
  nodes.set(data.id, [data.successor, data.d]);
  // console.log(nodes);
  update();
});

function update() {
  let n = [],
    p: [number, number][] = [];
  for (let node of nodes.keys()) n.push(parseInt(`0x${node.substring(0, 3)}`));
  for (const [key, val] of nodes) {
    for (let v of val)
      p.push([
        parseInt(`0x${key.substring(0, 3)}`),
        parseInt(`0x${v.substring(0, 3)}`),
      ]);
  }

  body.innerHTML = ``;
  renderer.set_peers(n).set_path(p).render();
  // console.log(n, nodes.keys(), nodes);
  console.log(p);
}
