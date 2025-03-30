const std = @import("std");
const c = @cImport(@cInclude("client.h"));

pub const Request = struct {
    url: [*c]u8,
    key: [*c]u8,
    route: [*c]u8,
    method: [*c]u8,
    input: [*c]u8,

    pub fn init(url: [*c]u8, key: [*c]u8, route: [*c]u8, method: [*c]u8, input: [*c]u8) Request {
        return Request{
            .url = url,
            .key = key,
            .route = route,
            .method = method,
            .input = input,
        };
    }

    pub fn MakeRequest(self: Request) []u8 {
        // Call the Go function, passing null for the input map
        const result = c.Request(self.url, self.key, self.route, self.method, self.input);
        defer std.c.free(@ptrCast(result)); // Free the C string when done
        // Convert the C string to a Zig string
        const result_str = std.mem.span(result);
        return result_str;
    }
};
