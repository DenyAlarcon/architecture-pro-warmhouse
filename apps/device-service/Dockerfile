# Build stage
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY DeviceService.csproj ./
RUN dotnet restore "DeviceService.csproj"
COPY . .
RUN dotnet publish "DeviceService.csproj" -c Release -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app/publish ./
EXPOSE 8082
ENV ASPNETCORE_URLS=http://+:8082
ENTRYPOINT ["dotnet", "DeviceService.dll"]
