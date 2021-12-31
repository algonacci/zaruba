from typing import Any, List, Mapping
from schemas.ztplAppCrudEntity import ZtplAppCrudEntity, ZtplAppCrudEntityData
from repos.ztplAppCrudEntity import ZtplAppCrudEntityRepo

class ZtplAppCrudEntityModel():

    def __init__(self, ztpl_app_crud_entity_repo: ZtplAppCrudEntityRepo):
        self.ztpl_app_crud_entity_repo = ztpl_app_crud_entity_repo

    def find(self, keyword: str, limit: int, offset: int) -> List[ZtplAppCrudEntity]:
        return self.ztpl_app_crud_entity_repo.find(keyword, limit, offset)

    def find_by_id(self, id: str) -> ZtplAppCrudEntity:
        return self.ztpl_app_crud_entity_repo.find_by_id(id)

    def insert(self, ztpl_app_crud_entity_data: ZtplAppCrudEntityData) -> ZtplAppCrudEntity:
        return self.ztpl_app_crud_entity_repo.insert(ztpl_app_crud_entity_data)

    def update(self, id: str, ztpl_app_crud_entity_data: ZtplAppCrudEntityData) -> ZtplAppCrudEntity:
        return self.ztpl_app_crud_entity_repo.update(id, ztpl_app_crud_entity_data)

    def delete(self, id: str) -> ZtplAppCrudEntity:
        return self.ztpl_app_crud_entity_repo.delete(id)