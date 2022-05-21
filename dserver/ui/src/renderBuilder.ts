import { DHTChord } from "./dht_chord_d3";

export class RenderBuilder {
  public chord: DHTChord;
  private path: [number, number][];

  constructor(
    root: string,
    radius: number,
    m: number,
    peers?: number[],
    path?: [number, number][]
  ) {
    this.chord = new DHTChord(root, radius, radius, m, peers || []);
    this.path = path || [];
  }

  set_peers(peers: number[]) {
    this.chord.peers = peers;

    return this;
  }

  set_path(path: [number, number][]) {
    this.path = path;
    return this;
  }

  render() {
    console.log(this.path);
    this.chord.build_chord();
    this.chord.draw_key_path(this.path);
    return this;
  }
}
