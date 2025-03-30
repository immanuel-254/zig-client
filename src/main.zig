const std = @import("std");
const client = @import("client.zig");

pub fn main() !void {
    // Define slices explicitly, ensuring they're valid
    const url: [*c]u8 = @constCast("http://localhost:3000".ptr);
    const key: [*c]u8 = @constCast("".ptr);
    const route: [*c]u8 = @constCast("/".ptr);
    const method: [*c]u8 = @constCast("GET".ptr);
    const input: [*c]u8 = @constCast("".ptr);

    var request = client.Request.init(url, key, route, method, input);

    // Get the response
    const result_str = request.MakeRequest();
    std.debug.print("Response: {s}\n", .{result_str});
}
