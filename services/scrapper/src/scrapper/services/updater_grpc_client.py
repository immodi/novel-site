from __future__ import print_function
import logging
import grpc
from scrapper.grpc import updater_pb2, updater_pb2_grpc


def run():
    # Connect to your gRPC server
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = updater_pb2_grpc.UpdaterServiceStub(channel)

        # Create the request
        request = updater_pb2.UpdateRequest(interval_hours=1)  # run every 1 hour

        print("Starting Updater stream...")

        try:
            # Stream server responses (since UpdateNovels yields multiple responses)
            for response in stub.UpdateNovels(request):
                print(f"📡 {response.message}")
        except grpc.RpcError as e:
            print(f"gRPC error: {e.code()} - {e.details()}")


if __name__ == "__main__":
    logging.basicConfig()
    run()
