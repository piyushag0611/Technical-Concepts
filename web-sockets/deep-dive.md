# WebSockets — Deep Dive

## 1. What is a WebSocket?

A WebSocket is a **persistent, full-duplex, bidirectional** communication channel over a single TCP connection. Once established, either side (client or server) can send messages at any time without the other side needing to request them first.

Compare to HTTP:
| | HTTP | WebSocket |
|---|---|---|
| Connection | New TCP connection per request | Single persistent TCP connection |
| Direction | Client → Server (request/response) | Both directions, anytime |
| Overhead | Headers on every message | Minimal framing after handshake |
| Use case | Fetch data | Real-time streaming, chat, games |

---

## 2. The Handshake — How a WebSocket is established

WebSockets **start as HTTP**, then upgrade. This is crucial to understand:

**Step 1 — Client sends HTTP upgrade request:**
```
GET /chat HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

**Step 2 — Server responds with 101 Switching Protocols:**
```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

The `Sec-WebSocket-Accept` is computed as:
```
Base64(SHA1(Sec-WebSocket-Key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
```
This magic GUID prevents cache poisoning attacks. After this, the TCP connection is no longer HTTP — it's a WebSocket connection.

---

## 3. The Wire Protocol — Frames

WebSocket messages are sent as **frames** (not raw TCP bytes):

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |                               |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+-------------------------------+
```

Key fields:
- **FIN bit** — is this the final fragment of a message?
- **Opcode** — `0x1` = text, `0x2` = binary, `0x8` = close, `0x9` = ping, `0xA` = pong
- **MASK bit** — client→server frames MUST be masked (prevents cache poisoning on proxies); server→client frames are NOT masked
- **Payload length** — 7-bit, or 16-bit, or 64-bit extended

---

## 4. Connection Lifecycle

```
Client                          Server
  |                               |
  |--- HTTP Upgrade Request ----→ |
  |← 101 Switching Protocols ---- |
  |                               |
  |←——— Full duplex channel ————→ |
  |                               |
  |--- TEXT frame ("hello") ----→ |
  |←--- TEXT frame ("world") ---- |
  |                               |
  |--- PING frame ---------------→|
  |←--- PONG frame -------------- |
  |                               |
  |--- CLOSE frame (1000) ------→ |
  |←--- CLOSE frame (1000) ------ |
  |                               |
  [TCP connection torn down]
```

Close codes:
- `1000` — Normal closure
- `1001` — Going away (page navigation)
- `1006` — Abnormal closure (no close frame — connection dropped)
- `1011` — Server error

---

## 5. Key Properties to Know Cold

**Full-duplex**: Client and server send independently — no polling, no waiting.

**Single TCP connection**: One port (80/443), one handshake, persistent.

**ws:// vs wss://**: wss is WebSocket over TLS — same as http vs https. Always use wss in production.

**No CORS for WebSockets**: The `Origin` header is sent but servers must validate it themselves — there's no browser-enforced CORS for WebSockets.

**Heartbeats (Ping/Pong)**: The protocol has built-in ping/pong frames. Used to detect dead connections and keep NAT/proxies from closing idle connections.

---

## 6. Scaling WebSockets — The Hard Part

HTTP is stateless — any server handles any request. WebSockets are **stateful** — a client is connected to one specific server process.

**Problem**: If you have 3 servers and Client A is on Server 1, Client B is on Server 2, and Server 1 wants to push a message to Client B — how?

**Solutions**:

1. **Pub/Sub broker (Redis, Kafka)** — All servers subscribe to a shared channel. When Server 1 needs to reach Client B, it publishes to the broker. Server 2 picks it up and pushes to Client B. This is the standard pattern.

2. **Sticky sessions** — Load balancer always routes a client to the same server (via IP hash or cookie). Simpler but creates hotspots and doesn't handle server failures well.

3. **Consistent hashing** — Clients are deterministically assigned to servers by ID. Better than sticky sessions but still stateful.

---

## 7. WebSockets vs. Alternatives

| Technology | Direction | Latency | Use Case |
|---|---|---|---|
| HTTP polling | Client-initiated | High | Simple, legacy |
| HTTP long polling | Client-initiated, held | Medium | When WS not available |
| SSE (Server-Sent Events) | Server → Client only | Low | Live feeds, notifications |
| WebSockets | Full-duplex | Low | Chat, games, collaboration |
| WebRTC | Peer-to-peer | Very low | Video, audio, P2P data |
| gRPC streaming | Full-duplex (HTTP/2) | Low | Service-to-service |

**When NOT to use WebSockets:**
- The client only needs to receive updates (use SSE — simpler)
- Infrequent updates with caching needs (use HTTP)
- You need request/response semantics with retries (HTTP with backoff)

---

## 8. Common Interview Questions

**Q: How does the WebSocket handshake prevent cross-protocol attacks?**
The Sec-WebSocket-Key/Accept exchange ensures the server is a real WebSocket server and not an HTTP cache that stored the upgrade response.

**Q: Why are client frames masked but server frames aren't?**
Masking prevents malicious scripts from sending predictable byte sequences that could poison transparent HTTP proxy caches. Servers are trusted infrastructure so masking isn't needed.

**Q: What happens when the network drops mid-connection?**
There's no close frame — the client and server don't know until a send/receive times out or a ping goes unanswered. This is close code 1006. Apps must implement heartbeat + reconnect logic themselves.

**Q: Can WebSockets work through HTTP proxies?**
Mostly yes, if the proxy supports HTTP/1.1 CONNECT tunneling. Strict proxies may block the Upgrade header — this is why wss:// (TLS) is preferred: traffic is encrypted so the proxy can't inspect it and tends to pass it through.

**Q: How many WebSocket connections can a server handle?**
Each connection is a file descriptor (socket). Linux default is ~1024 FDs per process, tunable to 1M+. With async I/O (Node.js, Go, Netty), a single server can handle 100k+ simultaneous connections. The bottleneck is usually memory (~10-50KB per connection) and the pub/sub backplane.

---

## 9. Infrastructure Cost — Why Every Layer Is Affected

The root cause: **WebSockets are persistent, stateful TCP connections**. Every component between client and server was designed for short-lived HTTP cycles. Holding a connection open for hours changes the operational model of each component.

### Load Balancer
- **HTTP model**: LB forwards request → server responds → connection closes. LB is stateless.
- **WebSocket problem**: LB must hold two TCP connections open (client↔LB, LB↔server) for hours. Standard round-robin breaks — a second request from the same client would land on a different server that has no WebSocket state.
- **Cost**: Requires sticky sessions (session affinity). On AWS, ALB charges per active connection-hour via LCU pricing — 10,000 idle connections costs real money even with no data flowing.

### Reverse Proxy / API Gateway (Nginx, Kong, Envoy)
- Must hold the TCP connection open and proxy frames in both directions indefinitely.
- HTTP-level timeouts will silently kill idle connections without explicit config.
- Nginx example — without this, connections drop after 60s of silence:
```nginx
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
proxy_read_timeout 3600s;
```
Every proxy in the chain needs this treatment. Miss one, and connections silently drop.

### CDN
- CDNs cache HTTP responses at edge nodes. WebSocket traffic cannot be cached.
- CDN must **tunnel** connections all the way back to origin — you lose edge latency benefits.
- Not all CDN tiers support WebSockets (Cloudflare requires Pro+, CloudFront needs configuration).
- Charged for connection time + data transfer, not just requests.

### Application Server — The Thread Model Problem
Traditional servers (Gunicorn, Puma, Tomcat) use **thread-per-request**: allocate a thread → handle request → free thread. With WebSockets, the thread is held for the connection's entire lifetime.

- 10,000 simultaneous connections = 10,000 threads = 10–80GB RAM just for thread stacks
- **Fix**: Use an async/event-loop server — Node.js, Go, Netty, Tornado, Actix. These handle thousands of connections on a small number of OS threads via non-blocking I/O. This often means re-platforming parts of the stack.

### Horizontal Scaling — The Pub/Sub Tax
Once you need more than one server, statefulness forces a message backplane:
```
Client A → Server 1 → publishes → Redis → delivers → Server 2 → Client B
```
- New infrastructure component to deploy, monitor, and scale (Redis Cluster)
- Every message now traverses 3 hops instead of 1
- New failure mode: Redis down = message delivery stops

### Kubernetes / Container Orchestration
- Rolling deploys become painful — Kubernetes sends `SIGTERM` but clients holding open connections may not reconnect gracefully, causing message loss.
- Need a graceful drain period and client-side reconnect logic.
- Can't auto-scale down by killing pods — must drain connections first.

### Firewalls and NAT
- Corporate firewalls and NAT gateways have idle connection timeout policies (typically 30s–5min).
- A silent WebSocket gets killed at the firewall; neither side knows until the next message fails.
- **This is why heartbeat ping/pong frames are mandatory in production** — they keep the connection alive through every intermediate device.

### Cost Summary

| Component | HTTP Cost | WebSocket Extra Cost |
|---|---|---|
| Load balancer | Per request | Per active connection-hour |
| Reverse proxy | Stateless | Must hold open + tune timeouts |
| CDN | Cheap (cached) | Full origin tunnel, no caching |
| App server | Thread freed after response | Thread/FD held for connection lifetime |
| Horizontal scale | None (stateless) | Pub/sub backplane required |
| Deploys | Simple rolling | Graceful drain required |
| Firewalls/NAT | Stateless | Heartbeat required to prevent silent drops |

### The Key Interview Insight
HTTP's statelessness is a feature that makes the entire stack simple to scale. WebSockets trade that simplicity for low-latency bidirectionality, and the cost shows up at every layer. This is why many teams reach for managed WebSocket services (AWS API Gateway WebSocket, Pusher, Ably) — they offload all this complexity at the cost of vendor lock-in.

---

## 10. Real-World Example — Does Netflix Use WebSockets?

**No — and the reasoning is a strong interview answer.**

Netflix's communication patterns don't call for full-duplex:

- **Video streaming** — Uses **DASH** (Dynamic Adaptive Streaming over HTTP) or **HLS**. The client requests the next chunk via HTTP range requests. The *client* drives the flow, not the server. HTTP is perfect.
- **User actions** (play, pause, seek, rate) — Fire-and-forget HTTP requests. No persistence needed.
- **"Are you still watching?"** — A client-side timer that fires a single HTTP request.
- **"New episode" notifications** — Low-frequency, can be handled by polling or SSE.

There is no scenario in Netflix's core product where the server needs to spontaneously push time-sensitive data to a client. No full-duplex = no need for WebSockets.

Netflix did build **Zuul 2**, an async non-blocking API gateway that supports WebSocket proxying — but this is for internal tooling and lower-traffic interactive features, not the main streaming path.

### The Communication Pattern Decision Table

| Pattern | Right tool |
|---|---|
| Client fetches data on demand | HTTP |
| Server streams updates to client only | SSE |
| Both sides send messages freely, anytime | WebSockets |
| Client drives a continuous stream | HTTP (chunked / range requests) |

Netflix lives in rows 1 and 4. Knowing *when not* to use WebSockets is as important as knowing how they work.

---

## 11. Group Video Calls — WebRTC vs WebSockets (WhatsApp Example)

**Short answer**: Group video calls use **WebRTC + SFU**. WebSockets handle signaling to set the call up. WebRTC carries the actual media. They work together, not as alternatives.

### Why Pure P2P WebRTC Breaks for Groups

In a 1:1 call, WebRTC P2P is ideal — media flows directly between peers. In a group call with N participants using pure **mesh topology** (everyone connects to everyone):
- Each participant uploads N-1 streams and downloads N-1 streams
- Bandwidth and CPU grow as **O(N²)**
- For an 8-person call, each person uploads 7 video streams — breaks down on mobile

### The SFU Architecture

Modern group video calls use a **Selective Forwarding Unit (SFU)** — a media server that routes (but does not mix/transcode) streams:

```
         ┌─────────────────────────────┐
         │            SFU              │
         │   (routes, not transcodes)  │
         └──┬──────┬──────┬──────┬────┘
            │      │      │      │
          WebRTC WebRTC WebRTC WebRTC
            │      │      │      │
          Alice   Bob   Carol   Dave
```

Each participant uploads 1 stream to the SFU and downloads N-1 streams from it. The SFU forwards raw RTP packets without decoding — low CPU, low latency. Still WebRTC end-to-end.

The alternative is an **MCU (Multipoint Control Unit)** which mixes all streams into one composite on the server — lower bandwidth for clients but massive server CPU and higher latency. Largely obsolete.

### Where WebSockets Fit In

**WebRTC deliberately does not define a signaling protocol.** It handles media transport but leaves call setup to you. Before two WebRTC peers can exchange media they need to:
1. Exchange **SDP offers/answers** — codec/resolution/port descriptions
2. Exchange **ICE candidates** — candidate network paths to try
3. Signal call lifecycle — ringing, accept, reject, hang up

This signaling uses **WebSockets**. Once signaling completes, media flows over WebRTC (UDP/SRTP) — the WebSocket is no longer in the hot path.

```
Alice            WhatsApp Server              Bob
  |-- SDP offer (WS) ------>|-- SDP offer (WS) ----->|
  |<-- SDP answer (WS) -----|<-- SDP answer (WS) -----|
  |-- ICE candidates (WS) ->|-- ICE candidates (WS)->|
  |<============ WebRTC media (UDP/SRTP) ============>|
```

### What WhatsApp Specifically Does

- **1:1 calls**: WebRTC P2P when possible, TURN relay when NAT/firewall blocks direct path
- **Group calls**: WebRTC + SFU (Meta runs their own SFU infrastructure)
- **Signaling**: WhatsApp's persistent connection (custom protocol, functionally equivalent to WebSockets)
- **Encryption**: Signal protocol applied on top of WebRTC media (end-to-end encrypted)

### The Mental Model

| Concern | Protocol |
|---|---|
| Real-time media (audio/video) | WebRTC (UDP/SRTP) |
| Call setup / signaling | WebSockets |
| Group call routing | SFU in the media path |
| Text chat during a call | WebSockets |

WebRTC and WebSockets are **complementary**. WebRTC solves low-latency media; WebSockets solve signaling and messaging. Almost every production video call system uses both.

### SFU Load — Packet Copies vs. Video Decoding

The SFU does send one copy of each packet per subscriber, but the cost is low because it operates at the **RTP packet level** — it never decodes video. An RTP packet is 1,200–1,400 bytes. The SFU just does:
1. Receive RTP packet from Alice
2. Look up routing table: who subscribed to Alice?
3. Copy the packet buffer to each subscriber's socket

This is what a **network switch** does, not what a **video codec** does. No H.264 decode, no frame manipulation, no re-encode. `memcpy` on a 1,400-byte buffer N times is cheap.

Compare to **MCU**: decodes every stream → composites frames → re-encodes. That's where the real CPU cost is. SFU entirely avoids it.

**The SFU's real bottleneck is outbound bandwidth, not CPU:**
- 100 participants × 1 Mbps each = 100 Mbps in
- Each participant receives 99 streams = ~10 Gbps out

**How SFUs manage this:**

- **Simulcast**: Senders encode multiple quality levels (1080p, 360p, 180p) simultaneously. SFU forwards only the quality each receiver's bandwidth supports — no transcoding, just picking a pre-encoded layer.
- **SVC (Scalable Video Coding)**: Video encoded in layers (base + enhancement). SFU drops enhancement layers for bandwidth-constrained receivers without re-encoding.
- **Active speaker detection**: SFU monitors audio levels via RTP header extensions, forwards active speaker at full quality, reduces or pauses others. This is why the speaking person's tile gets large in Meet/Zoom.
- **Bandwidth estimation (REMB/TWCC)**: Receivers send RTCP feedback reporting available bandwidth. SFU adjusts forwarding in real time.

| | MCU | SFU |
|---|---|---|
| Decodes video? | Yes (every stream) | No |
| Re-encodes? | Yes (composite output) | No |
| CPU cost | Very high | Low |
| Bandwidth cost | Low (1 stream out) | High (N-1 streams out) |
| Latency | High | Low |
| Mitigation | N/A | Simulcast, SVC, active speaker |

The SFU trades CPU for bandwidth. Bandwidth scales horizontally and cheaply; CPU does not.
