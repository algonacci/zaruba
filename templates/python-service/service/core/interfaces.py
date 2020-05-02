from abc import ABC, abstractmethod
from flask import Flask
from logging import Logger
from typing import List, Callable

SetupComponent = Callable[[], None]


class App(ABC):

    @abstractmethod
    def logger(self) -> Logger:
        pass

    @abstractmethod
    def router(self) -> Flask:
        pass

    @abstractmethod
    def liveness(self) -> bool:
        pass

    @abstractmethod
    def readiness(self) -> bool:
        pass

    @abstractmethod
    def set_liveness(self, liveness: bool) -> None:
        pass

    @abstractmethod
    def set_readiness(self, readiness: bool) -> None:
        pass

    @abstractmethod
    def setup(self, setupComponents: List[SetupComponent]) -> None:
        pass

    @abstractmethod
    def run(self) -> None:
        pass
