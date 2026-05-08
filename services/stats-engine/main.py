import os
from concurrent.futures.thread import ThreadPoolExecutor
import grpc
from server.stats_engine_server import StatsEngineServer
import stats_engine.v1.stats_engine_pb2_grpc as ses_grpc
import stats_engine.v1.stats_engine_pb2 as ses
from dotenv import load_dotenv
from grpc_reflection.v1alpha import reflection
import pandas as pd
import spotify_confidence as conf
import warnings
warnings.simplefilter(action='ignore', category=FutureWarning)



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
    # df = pd.DataFrame({'variation_name': ["control", "treatment"],
    #                    'success': [3000, 3500],
    #                    'total': [5000, 5000],
    #                    })
    # df.head()
    #
    # ztest = conf.ZTest(data_frame=df,
    #                    numerator_column='success',
    #                    numerator_sum_squares_column='success',
    #                    denominator_column='total',
    #                    interval_size=0.95,
    #                    correction_method='bonferroni',
    #                    categorical_group_columns = 'variation_name',
    #                    )
    #
    # print(ztest.summary().to_string())
    # print(ztest.difference("treatment", "control").to_string())
    load_dotenv("../../infrastructure/.env")
    serve()


