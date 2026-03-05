#!/bin/bash

# Source environment variables
if [ -f .env ]; then
  set -a
  source .env
  set +a
else
  echo "Error: .env file not found"
  exit 1
fi

# Generate timestamp tag in format YYYYMMDDHHMMSS
TAG=$(date '+%Y%m%d%H%M%S')
FULL_TAG="${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${TAG}"

echo "Building and pushing Docker images with tag: ${TAG}"
echo "Registry: ${DOCKER_REGISTRY}"
echo "Image Name: ${DOCKER_IMAGE_NAME}"

# Build broker image
echo "Building broker image..."
docker build -f docker/broker.Dockerfile \
  -t "${FULL_TAG}" \
  -t "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:broker-latest" \
  .

if [ $? -ne 0 ]; then
  echo "Error building broker image"
  exit 1
fi

# Push broker image
echo "Pushing broker image..."
docker push "${FULL_TAG}"
docker push "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:broker-latest"

# Build admin image
echo "Building admin image..."
docker build -f docker/admin.Dockerfile \
  -t "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${TAG}-admin" \
  -t "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:admin-latest" \
  .

if [ $? -ne 0 ]; then
  echo "Error building admin image"
  exit 1
fi

# Push admin image
echo "Pushing admin image..."
docker push "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${TAG}-admin"
docker push "${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:admin-latest"

echo "Successfully built and pushed images:"
echo "  Broker: ${FULL_TAG}"
echo "  Admin:  ${DOCKER_REGISTRY}/${DOCKER_IMAGE_NAME}:${TAG}-admin"
