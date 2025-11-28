import time
import grpc
import argparse
from scrapper import config
from scrapper.cache.db_cache import NovelDataCache
from scrapper.grpc_services import updater_pb2, updater_pb2_grpc


def restart_update_service():
    """
    Safely restarts the updater gRPC service.
    """
    print("[UpdateCheck] Restarting gRPC updater service...")

    try:
        with grpc.insecure_channel("localhost:50051") as channel:
            stub = updater_pb2_grpc.UpdaterServiceStub(channel)

            try:
                stub.StopUpdate(updater_pb2.Empty())
            except Exception:
                pass  # ignore stop failure

            start_req = updater_pb2.UpdateRequest(interval_hours=3)
            resp = stub.StartUpdate(start_req)
            print(f"[Updater] {resp.message}")

    except Exception as e:
        print(f"[UpdateCheck] ERROR restarting updater service: {e}")


def update_check(interval_hours_minutes: int):
    """
    Performs one check and restart if needed.
    """
    print("[UpdateCheck] Running update check...")

    try:
        cache = NovelDataCache(config.CACHE_DB_PATH)
    except Exception as e:
        print(f"[UpdateCheck] ERROR opening cache: {e}")
        return

    try:
        try:
            last_update = cache.get_last_update()
        except Exception as e:
            print(f"[UpdateCheck] ERROR reading last_update: {e}")
            last_update = None

        now_minutes = int(time.time() // 60)

        should_restart = (
            last_update is None or (now_minutes - last_update) >= interval_hours_minutes
        )

        if should_restart:
            print("[UpdateCheck] Triggering updater restart...")
            restart_update_service()

            try:
                cache.set_last_update(now_minutes)
            except Exception as e:
                print(f"[UpdateCheck] ERROR writing last_update: {e}")

        else:
            if last_update is not None:
                diff = now_minutes - last_update
                print(f"[UpdateCheck] Skip — last update was {diff} minutes ago.")

    finally:
        try:
            cache.close()
        except Exception:
            pass


def run_update_monitor():
    """
    Runs update checks forever
    """
    parser = argparse.ArgumentParser(
        prog="updater-check", description="Novel updater check CLI"
    )
    parser.add_argument(
        "interval_hours",
        type=int,
        default=6,
        help="Interval in hours between updates (default: 6)",
    )
    args = parser.parse_args()
    interval_hours = args.interval_hours

    print(f"Updater started. Interval hours: {interval_hours}")

    while True:
        update_check(interval_hours * 60)
        print(f"[UpdateCheck] Sleeping for {interval_hours} hours...\n")
        time.sleep(interval_hours * 60 * 60)


if __name__ == "__main__":
    run_update_monitor()
