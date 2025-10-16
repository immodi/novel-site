import threading
import queue
import time
from scrapper.grpc import updater_pb2


class UpdateManager:
    def __init__(self):
        self.running = False
        self.thread = None
        self.listeners = []
        self.interval_hours = 3

    def start(self, update_fn, interval_hours: int):
        if self.running:
            return "Updater already running"
        self.running = True
        self.interval_hours = interval_hours
        self.thread = threading.Thread(
            target=self._run_loop, args=(update_fn,), daemon=True
        )
        self.thread.start()
        return "Updater started."

    def stop(self):
        if not self.running:
            return "Updater is not running"
        self.broadcast("Updater service stopped.")
        self.running = False
        return "Updater stopped."

    def subscribe(self):
        q = queue.Queue()
        q.put(
            updater_pb2.UpdateResponse(  # type: ignore
                message=f"🚀 UPDATER SERVICE started with an interval of {self.interval_hours} hour(s)"
            )
        )
        self.listeners.append(q)
        return q

    def unsubscribe(self, q):
        if q in self.listeners:
            self.listeners.remove(q)

    def broadcast(self, msg: str):
        for q in self.listeners:
            q.put(updater_pb2.UpdateResponse(message=msg))  # type: ignore

    def _run_loop(self, update_fn):
        while self.running:
            try:
                self.broadcast("🔁 Starting update cycle...")
                update_fn(self.broadcast)
                self.broadcast("✅ Update cycle complete.")
            except Exception as e:
                self.broadcast(f"❌ Error in update loop: {e}")
            time.sleep(self.interval_hours * 3600)
