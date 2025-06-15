#!/usr/bin/env bash
set -euo pipefail

# Load environment variables
if [ -f .env ]; then
  # shellcheck disable=SC1091
  source .env
else
  echo ".env file not found!" >&2
  exit 1
fi

# Usage: ./publish.sh [tag]
TAG=${1:-${DH_TAG}}
IMAGE="${DH_USER}/${DH_REPO}:${TAG}"

echo "🔨 Building ${IMAGE}..."
docker build -t "${IMAGE}" .

echo "🚀 Pushing to Docker Hub..."
docker push "${IMAGE}"

echo "✅ Pushed ${IMAGE}"