from .interfaces import App, SetupComponent
from flask import Flask
from logging import Logger, getLogger
from typing import List
import asyncio


class MainApp(App):

    def __init__(self, http_port: int, global_rmq_connection_string: str, local_rmq_connection_string: str):
        self._http_port: int = http_port
        self._app: Flask = Flask(__name__.split(".")[0])
        self._logger: Logger = getLogger(__name__.split(".")[0])
        self._liveness: bool = False
        self._readiness: bool = False

    def logger(self) -> Logger:
        return self._logger

    def router(self) -> Flask:
        return self._app

    def liveness(self) -> bool:
        return self._liveness

    def readiness(self) -> bool:
        return self._readiness

    def set_liveness(self, liveness: bool) -> None:
        self._liveness = liveness

    def set_readiness(self, readiness: bool) -> None:
        self._readiness = readiness

    def setup(self, setupComponents: List[SetupComponent]) -> None:
        for setupComponent in setupComponents:
            setupComponent()

    def run(self) -> None:
        try:
            with self._app.app_context():
                print("hi")
            self._app.run("0.0.0.0", self._http_port)
        except Exception as e:
            self._logger.error(e)