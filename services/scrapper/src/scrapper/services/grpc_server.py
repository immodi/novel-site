from scrapper.grpc import scrapper_pb2_grpc, updater_pb2_grpc
from scrapper.services.endpoints.updater_grpc_endpoint import UpdaterService
from scrapper.services.endpoints.scrapper_grpc_endpoint import ScrapperService
from concurrent import futures
import logging
import grpc


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    scrapper_pb2_grpc.add_ScrapperServiceServicer_to_server(ScrapperService(), server)
    updater_pb2_grpc.add_UpdaterServiceServicer_to_server(UpdaterService(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
