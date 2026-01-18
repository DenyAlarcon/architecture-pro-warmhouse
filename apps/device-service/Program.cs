using System.Collections.Concurrent;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var devices = new ConcurrentDictionary<string, Device>();

app.MapPost("/devices", (DeviceCreate request) =>
{
    var deviceId = string.IsNullOrWhiteSpace(request.DeviceId)
        ? Guid.NewGuid().ToString("D")
        : request.DeviceId;

    var device = new Device(
        deviceId,
        request.Name ?? string.Empty,
        request.Type ?? string.Empty,
        request.Location ?? string.Empty,
        request.Status ?? "inactive",
        request.Unit
    );

    devices[device.Id] = device;
    return Results.Created($"/devices/{device.Id}", device);
});

app.MapGet("/devices/{deviceId}", (string deviceId) =>
{
    return devices.TryGetValue(deviceId, out var device)
        ? Results.Ok(device)
        : Results.NotFound();
});

app.MapPatch("/devices/{deviceId}/status", (string deviceId, DeviceStatusUpdate request) =>
{
    if (!devices.TryGetValue(deviceId, out var device))
    {
        return Results.NotFound();
    }

    var updated = device with { Status = request.Status };
    devices[deviceId] = updated;
    return Results.Ok(updated);
});

app.Run();

record Device(
    string Id,
    string Name,
    string Type,
    string Location,
    string Status,
    string? Unit
);

record DeviceCreate(
    string? DeviceId,
    string? Name,
    string? Type,
    string? Location,
    string? Status,
    string? Unit
);

record DeviceStatusUpdate(string Status);
