from typing import Optional
from schemas.ztplAppCrudEntity import ZtplAppCrudEntity, ZtplAppCrudEntityData, ZtplAppCrudEntityResult
from repos.ztplAppCrudEntity import ZtplAppCrudEntityRepo

class ZtplAppCrudEntityService():

    def __init__(self, ztpl_app_crud_entity_repo: ZtplAppCrudEntityRepo):
        self.ztpl_app_crud_entity_repo = ztpl_app_crud_entity_repo

    def find(self, keyword: str, limit: int, offset: int) -> ZtplAppCrudEntityResult:
        count = self.ztpl_app_crud_entity_repo.count(keyword)
        rows = self.ztpl_app_crud_entity_repo.find(keyword, limit, offset)
        return ZtplAppCrudEntityResult(count=count, rows=rows)

    def find_by_id(self, id: str) -> Optional[ZtplAppCrudEntity]:
        return self.ztpl_app_crud_entity_repo.find_by_id(id)

    def insert(self, ztpl_app_crud_entity_data: ZtplAppCrudEntityData) -> Optional[ZtplAppCrudEntity]:
        return self.ztpl_app_crud_entity_repo.insert(ztpl_app_crud_entity_data)

    def update(self, id: str, ztpl_app_crud_entity_data: ZtplAppCrudEntityData) -> Optional[ZtplAppCrudEntity]:
        return self.ztpl_app_crud_entity_repo.update(id, ztpl_app_crud_entity_data)

    def delete(self, id: str) -> Optional[ZtplAppCrudEntity]:
        return self.ztpl_app_crud_entity_repo.delete(id)