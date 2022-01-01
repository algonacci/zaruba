from typing import Any, List, Mapping
from helpers.transport import RPC
from schemas.ztplAppCrudEntity import ZtplAppCrudEntityData
from schemas.user import User
from repos.ztplAppCrudEntity import ZtplAppCrudEntityRepo
from ztplAppModuleName.ztplAppCrudEntityModel import ZtplAppCrudEntityModel

def register_ztpl_app_crud_entity_rpc(rpc: RPC, ztpl_app_crud_entity_repo: ZtplAppCrudEntityRepo):

    ztpl_app_crud_entity_model = ZtplAppCrudEntityModel(ztpl_app_crud_entity_repo)

    @rpc.handle('find_ztpl_app_crud_entity')
    def find_ztpl_app_crud_entity(keyword: str, limit: int, offset: int, current_user_data: Mapping[str, Any]) -> List[Mapping[str, Any]]:
        results = ztpl_app_crud_entity_model.find(keyword, limit, offset)
        return [result.dict() for result in results]

    @rpc.handle('find_ztpl_app_crud_entity_by_id')
    def find_ztpl_app_crud_entity_by_id(id: str, current_user_data: Mapping[str, Any]) -> Mapping[str, Any]:
        result = ztpl_app_crud_entity_model.find_by_id(id)
        return None if result is None else result.dict()

    @rpc.handle('insert_ztpl_app_crud_entity')
    def insert_ztpl_app_crud_entity(ztpl_app_crud_entity_data: Mapping[str, Any], current_user_data: Mapping[str, Any]) -> Mapping[str, Any]:
        result = ztpl_app_crud_entity_model.insert(ZtplAppCrudEntityData.parse_obj(ztpl_app_crud_entity_data))
        return None if result is None else result.dict()

    @rpc.handle('update_ztpl_app_crud_entity')
    def update_ztpl_app_crud_entity(id: str, ztpl_app_crud_entity_data: Mapping[str, Any], current_user_data: Mapping[str, Any]) -> Mapping[str, Any]:
        result = ztpl_app_crud_entity_model.update(id, ZtplAppCrudEntityData.parse_obj(ztpl_app_crud_entity_data))
        return None if result is None else result.dict()

    @rpc.handle('delete_ztpl_app_crud_entity')
    def delete_ztpl_app_crud_entity(id: str, current_user_data: Mapping[str, Any]) -> Mapping[str, Any]:
        result = ztpl_app_crud_entity_model.delete(id)
        return None if result is None else result.dict()

    print('Handle RPC for ztplAppModuleName.ZtplAppCrudEntity')