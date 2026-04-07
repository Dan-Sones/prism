import os
from concurrent.futures.thread import ThreadPoolExecutor
import grpc
from server.stats_engine_server import StatsEngineServer
import stats_engine.v1.stats_engine_pb2_grpc as ses_grpc
import stats_engine.v1.stats_engine_pb2 as ses
from dotenv import load_dotenv
from grpc_reflection.v1alpha import reflection


def serve():
    port = os.getenv("STATS_ENGINE_GRPC_SERVER_PORT")
    server = grpc.server(ThreadPoolExecutor(max_workers=10))
    ses_grpc.add_StatsEngineServicer_to_server(StatsEngineServer(), server)

    service_names = (
        ses.DESCRIPTOR.services_by_name['StatsEngine'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(service_names, server)

    print("Starting Stats Engine gRPC server on port {port}...".format(port=port))
    server.add_insecure_port("[::]:{port}".format(port=port))
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    load_dotenv("../../infrastructure/.env")
    serve()


