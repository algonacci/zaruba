from typing import Callable, List, Optional
from fastapi.security import OAuth2PasswordBearer, OAuth2
from fastapi import Depends, Cookie, HTTPException, status
from starlette.requests import Request
from auth.roleService import RoleService
from auth.userService import UserService
from auth.tokenService import TokenService
from schemas.user import User

import abc
import traceback

class AuthService(abc.ABC):

    @abc.abstractmethod
    def everyone(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        pass

    @abc.abstractmethod
    def is_authenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        pass

    @abc.abstractmethod
    def is_unauthenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        pass

    @abc.abstractmethod
    def is_authorized(self, permission: str, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        pass

class NoAuthService(AuthService):

    def __init__(self, user_service: UserService):
        self.user_service: str = user_service

    def everyone(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        return self.user_service.get_guest_user()

    def is_unauthenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        if throw_error:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail='Forbidden'
            )
        return None

    def is_authenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        return self.user_service.get_guest_user()

    def is_authorized(self, permission: str, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        return self.user_service.get_guest_user()

class TokenOAuth2AuthService(AuthService):

    def __init__(self, user_service: UserService, token_service: TokenService, oauth2_scheme: OAuth2):
        self.user_service = user_service
        self.token_service = token_service
        self.oauth2_scheme = oauth2_scheme

    def _raise_error_or_return_none(self, throw_error: bool, status_code: int, detail: str) -> None:
        if not throw_error:
            return None
        raise HTTPException(
            status_code=status_code,
            detail=detail,
            headers={'WWW-Authenticate': 'Bearer'},
        )

    def _get_user_by_token(self, token: str) -> Optional[User]:
        try:
            return self.token_service.get_user_by_token(token)
        except:
            print(traceback.format_exc)
            return None

    def everyone(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        async def verify_everyone(bearer_token = Depends(self.oauth2_scheme), app_access_token=Cookie(default=None)) -> Optional[User]:
            if bearer_token is None and app_access_token is None:
                return None
            token = bearer_token if bearer_token is not None else app_access_token
            return self._get_user_by_token(token)
        return verify_everyone 

    def is_unauthenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        async def verify_is_unauthenticated(bearer_token = Depends(self.oauth2_scheme), app_access_token=Cookie(default=None)) -> Optional[User]:
            if bearer_token is None and app_access_token is None:
                return None
            token = bearer_token if bearer_token is not None else app_access_token
            current_user = self._get_user_by_token(token)
            if not current_user or not current_user.active:
                return None
            return self._raise_error_or_return_none(throw_error, status.HTTP_401_UNAUTHORIZED, 'Not authenticated')
        return verify_is_unauthenticated

    def is_authenticated(self, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        async def verify_is_authenticated(bearer_token = Depends(self.oauth2_scheme), app_access_token=Cookie(default=None)) -> Optional[User]:
            if bearer_token is None and app_access_token is None:
                return self._raise_error_or_return_none(throw_error, status.HTTP_401_UNAUTHORIZED, 'Not authenticated')
            token = bearer_token if bearer_token is not None else app_access_token
            current_user = self._get_user_by_token(token)
            if not current_user or not current_user.active:
                return self._raise_error_or_return_none(throw_error, status.HTTP_401_UNAUTHORIZED, 'Not authenticated')
            return current_user
        return verify_is_authenticated

    def is_authorized(self, permission: str, throw_error: bool = True) -> Callable[[Request], Optional[User]]:
        async def verify_is_authorized(current_user = Depends(self.is_authenticated(throw_error=throw_error))) -> Optional[User]:
            if not current_user:
                return self._raise_error_or_return_none(throw_error, status.HTTP_403_FORBIDDEN, 'Not authenticated')
            if self.user_service.is_authorized(current_user, permission):
                return current_user
            self._raise_error_or_return_none(throw_error, status.HTTP_403_FORBIDDEN, 'Unauthorized')
        return verify_is_authorized
   