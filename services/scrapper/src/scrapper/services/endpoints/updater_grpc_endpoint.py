from scrapper.services.modules.update_manager import UpdateManager
from scrapper.services.modules.update_runner import run_update
from scrapper.grpc_services import updater_pb2, updater_pb2_grpc
import queue

update_manager = UpdateManager()


class UpdaterService(updater_pb2_grpc.UpdaterServiceServicer):
    def StartUpdate(self, request, context):
        msg = update_manager.start(run_update, request.interval_hours)
        return updater_pb2.UpdateResponse(message=msg)  # type: ignore

    def StopUpdate(self, request, context):
        msg = update_manager.stop()
        return updater_pb2.UpdateResponse(message=msg)  # type: ignore

    def StreamUpdates(self, request, context):
        q = update_manager.subscribe()
        try:
            while context.is_active():
                try:
                    msg = q.get(timeout=1)
                    yield msg
                except queue.Empty:
                    continue
        finally:
            update_manager.unsubscribe(q)
