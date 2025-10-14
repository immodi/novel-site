import grpc
from scrapper.grpc import updater_pb2, updater_pb2_grpc


def stop_updater():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = updater_pb2_grpc.UpdaterServiceStub(channel)
        resp = stub.StopUpdate(updater_pb2.Empty())
        print(f"🛑 {resp.message}")


if __name__ == "__main__":
    stop_updater()
