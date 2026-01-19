using System.Collections.Concurrent;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var telemetryStore = new ConcurrentDictionary<string, List<TelemetryRecord>>();

app.MapPost("/telemetry", (TelemetryRecordCreate request) =>
{
    var record = new TelemetryRecord(
        request.DeviceId,
        request.Metric ?? "value",
        request.Value,
        request.Unit,
        request.RecordedAt ?? DateTime.UtcNow
    );

    var list = telemetryStore.GetOrAdd(record.DeviceId, _ => new List<TelemetryRecord>());
    lock (list)
    {
        list.Add(record);
    }

    return Results.Accepted();
});

app.MapGet("/telemetry/{deviceId}", (string deviceId, int? limit) =>
{
    if (!telemetryStore.TryGetValue(deviceId, out var list))
    {
        return Results.Ok(new TelemetryResponse(deviceId, new List<TelemetryRecord>()));
    }

    List<TelemetryRecord> snapshot;
    lock (list)
    {
        snapshot = list
            .OrderByDescending(x => x.RecordedAt)
            .Take(limit ?? 50)
            .ToList();
    }

    return Results.Ok(new TelemetryResponse(deviceId, snapshot));
});

app.Run();

record TelemetryRecord(
    string DeviceId,
    string Metric,
    double Value,
    string? Unit,
    DateTime RecordedAt
);

record TelemetryRecordCreate(
    string DeviceId,
    string? Metric,
    double Value,
    string? Unit,
    DateTime? RecordedAt
);

record TelemetryResponse(string DeviceId, List<TelemetryRecord> Items);
