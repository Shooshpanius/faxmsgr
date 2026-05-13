# =============================================================================
# Stage 1: сборка
# =============================================================================
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
ARG BUILD_CONFIGURATION=Release

# Установка Node.js 22.x для сборки SPA-проекта (Vite требует Node.js 20.19+ или 22.12+)
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl ca-certificates \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# Восстановление зависимостей .NET (кэш слоя)
COPY ["faxmsgr.Server/faxmsgr.Server.csproj", "faxmsgr.Server/"]
RUN dotnet restore "faxmsgr.Server/faxmsgr.Server.csproj"

# Установка зависимостей npm (кэш слоя)
COPY ["faxmsgr.client/package.json", "faxmsgr.client/package-lock.json*", "faxmsgr.client/"]
RUN cd faxmsgr.client && npm ci

# Копирование исходников и публикация
COPY . .
RUN dotnet publish "faxmsgr.Server/faxmsgr.Server.csproj" \
    -c $BUILD_CONFIGURATION \
    -o /app/publish \
    --no-restore

# =============================================================================
# Stage 2: runtime
# =============================================================================
FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime

# libgssapi-krb5-2 требуется Npgsql для подключения к PostgreSQL через GSSAPI
RUN apt-get update \
    && apt-get install -y --no-install-recommends libgssapi-krb5-2 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Открываем HTTP-порт (рекомендуемый порт ASP.NET Core в контейнерах)
EXPOSE 8080

COPY --from=build /app/publish .

ENTRYPOINT ["dotnet", "faxmsgr.Server.dll"]
