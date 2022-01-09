from os import access
from typing import Any, List, Mapping
from auth.authModel import AuthModel
from helpers.transport import MessageBus, RPC
from fastapi import Depends, FastAPI, HTTPException
from schemas.user import User, UserData

import traceback


def register_user_route(app: FastAPI, mb: MessageBus, rpc: RPC, auth_model: AuthModel):

    @app.get('/users/', response_model=List[User])
    def find_user(keyword: str='', limit: int=100, offset: int=0, current_user = Depends(auth_model.has_any_permissions( 'user:read'))) -> List[User]:
        results = []
        try:
            results = rpc.call('find_user', keyword, limit, offset, current_user.dict())
        except:
            print(traceback.format_exc()) 
            raise HTTPException(status_code=500, detail='Internal Server Error')
        return [User.parse_obj(result) for result in results]

    @app.get('/users/{id}', response_model=User)
    def find_user_by_id(id: str, current_user = Depends(auth_model.has_any_permissions( 'user:read'))) -> User:
        result = None
        try:
            result = rpc.call('find_user_by_id', id, current_user.dict())
        except:
            print(traceback.format_exc()) 
            raise HTTPException(status_code=500, detail='Internal Server Error')
        if result is None:
            raise HTTPException(status_code=404, detail='Not Found')
        return User.parse_obj(result)

    @app.post('/users/', response_model=User)
    def insert_user(data: UserData, current_user = Depends(auth_model.has_any_permissions( 'user:create'))) -> User:
        result = None
        try:
            result = rpc.call('insert_user', data.dict(), current_user.dict())
        except:
            print(traceback.format_exc()) 
            raise HTTPException(status_code=500, detail='Internal Server Error')
        if result is None:
            raise HTTPException(status_code=404, detail='Not Found')
        return User.parse_obj(result)

    @app.put('/users/{id}', response_model=User)
    def update_user(id: str, data: UserData, current_user = Depends(auth_model.has_any_permissions( 'user:update'))) -> User:
        result = None
        try:
            result = rpc.call('update_user', id, data.dict(), current_user.dict())
        except:
            print(traceback.format_exc()) 
            raise HTTPException(status_code=500, detail='Internal Server Error')
        if result is None:
            raise HTTPException(status_code=404, detail='Not Found')
        return User.parse_obj(result)

    @app.delete('/users/{id}')
    def delete_user(id: str, current_user = Depends(auth_model.has_any_permissions( 'user:delete'))) -> User:
        result = None
        try:
            result = rpc.call('delete_user', id)
        except:
            print(traceback.format_exc()) 
            raise HTTPException(status_code=500, detail='Internal Server Error')
        if result is None:
            raise HTTPException(status_code=404, detail='Not Found')
        return User.parse_obj(result)

    print('Handle HTTP routes for auth.User')