using System.Text;
using System.Net.WebSockets;

namespace WebSocketDemo;

public class WebSocketTests
{
    // Server starts and listens on a port.
    [Fact]
    public async Task ServerStartsAndListens()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);
        server.Stop();
        Assert.False(server.IsListening);
    }

    // Client can connect to the server.
    [Fact]
    public async Task ClientCanConnect()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        server.Stop();
        Assert.False(server.IsListening);
    }

    //On connection, client receives a welcome message from the server.
    [Fact]
    public async Task ClientReceivesWelcomeOnConnection()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        server.Stop();
        Assert.False(server.IsListening);
    }

    // Client sends a message, and receives the echoed message coorectly

    [Fact]
    public async Task ClientMessagesReceivesEcho()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        string clientMessage = "Hello";
        var echoMessage = await client.SendMessage(conn, clientMessage);

        Assert.Equal(clientMessage, echoMessage);
        server.Stop();
        Assert.False(server.IsListening);
    }

    //Server sends an independent message, and client receives the message correctly.
    [Fact]
    public async Task ServerMessagesClientRecieves()
    {
        string serverMessage = "Ping!";
        var server = new Server();
        await server.StartAsync(
            handler: async(ws) =>
            {
                var bytes = Encoding.UTF8.GetBytes(Server.WELCOME_MESSAGE);
                await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);

                bytes = Encoding.UTF8.GetBytes(serverMessage);
                await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
                
                var buffer = new byte[1024];
                await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
            }
        );
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        buffer = new byte[1024];
        result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(serverMessage, message);

        server.Stop();
        Assert.False(server.IsListening);
    }

    //Client can close the connection cleanly
    [Fact]
    public async Task ClientCloseConnection()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        WebSocketState state = await client.CloseConnection(conn);

        Assert.Equal(WebSocketState.Closed, state);

        server.Stop();
        Assert.False(server.IsListening);
    }

    //Server handles multiple sequential messages correctly.
    [Fact]
    public async Task ServerSequentialMessages()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        List<string> clientMessages = ["Hello", "how", "do", "you", "do"];
        List<string> receivedMessages = [];
        foreach (var clientMessage in clientMessages)
        {
            var echoMessage = await client.SendMessage(conn, clientMessage);
            receivedMessages.Add(echoMessage);
        }
        Assert.Equal(clientMessages, receivedMessages);
        server.Stop();
        Assert.False(server.IsListening);
    }

    //Client handles multiple sequential messages correctly.
    [Fact]
    public async Task ClientSequentialMessages()
    {
        List<string> serverMessages = ["I", "am", "doing", "well", "thank", "you"];
        var server = new Server();
        await server.StartAsync(
            handler: async(ws) =>
            {
                var bytes = Encoding.UTF8.GetBytes(Server.WELCOME_MESSAGE);
                await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);

                foreach (var message in serverMessages)
                {
                    bytes = Encoding.UTF8.GetBytes(message);
                    await ws.SendAsync(bytes, WebSocketMessageType.Text, true, CancellationToken.None);
                }
                var buffer = new byte[1024];
                await ws.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
            }
        );
        Assert.True(server.IsListening);

        var client = new Client();
        var conn = await client.ConnectAsync();

        Assert.Equal(WebSocketState.Open, conn.State);

        var buffer = new byte[1024];
        var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
        var message = Encoding.UTF8.GetString(buffer, 0, result.Count);
        Assert.Equal(Server.WELCOME_MESSAGE, message);

        var serverMessagesClient = await client.ReceiveServerMessages(conn);
        
        Assert.Equal(serverMessages, serverMessagesClient);
        server.Stop();
        Assert.False(server.IsListening);
    }

    //Server can handle multiple clients
    [Fact]
    public async Task HandleMultipleClients()
    {
        var server = new Server();
        await server.StartAsync();
        Assert.True(server.IsListening);

        var client1 = new Client();
        var client2 = new Client();
        var client3 = new Client();

        List<Client> clients = [client1, client2, client3];
        var tasks = clients.Select(client => client.ConnectAsync());
        var conns = await Task.WhenAll(tasks);

        foreach (var conn in conns)
        {
            var buffer = new byte[1024];
            var result = await conn.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);
            Assert.Equal(Server.WELCOME_MESSAGE, Encoding.UTF8.GetString(buffer, 0, result.Count));
        }

        List<string> clientMessages = ["Ping1", "Ping2", "Ping3"];

        var messageTasks = Enumerable.Range(0, 3).Select(index=> clients[index].SendMessage(conns[index], clientMessages[index]));
        var echoMessages = await Task.WhenAll(messageTasks);
        Assert.Equal(clientMessages, echoMessages);
        server.Stop();
        Assert.False(server.IsListening);
    }
}
