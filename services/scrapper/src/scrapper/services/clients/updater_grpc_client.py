from __future__ import print_function
import logging
import grpc
from scrapper.grpc import updater_pb2, updater_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = updater_pb2_grpc.UpdaterServiceStub(channel)

        # 1️⃣ Start the background updater
        start_req = updater_pb2.UpdateRequest(interval_hours=1)
        resp = stub.StartUpdate(start_req)
        print(f"{resp.message}")

        # 2️⃣ Listen to the update stream
        print("📡 Streaming updates...")
        try:
            for response in stub.StreamUpdates(updater_pb2.Empty()):
                print(f"{response.message}")
        except grpc.RpcError as e:
            print(f"Stream ended: {e.code()} - {e.details()}")

        # 3️⃣ Optionally stop it
        # stop_resp = stub.StopUpdate(updater_pb2.Empty())
        # print(f"{stop_resp.message}")


if __name__ == "__main__":
    logging.basicConfig()
    run()
