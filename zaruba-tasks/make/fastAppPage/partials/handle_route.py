
if enable_ui:
    @app.get('ztplAppUrl', response_class=HTMLResponse)
    async def get_ztpl_app_url(request: Request, context: MenuContext = Depends(menu_service.authenticate('ztplAppModuleName:ztplAppUrl'))) -> HTMLResponse:
        '''
        Handle (get) ztplAppUrl
        '''
        try:
            return templates.TemplateResponse('default_page.html', context={
                'request': request,
                'context': context,
                'content_path': 'ztplAppModuleName/ztpl_app_url.html'
            }, status_code=200)
        except:
            print(traceback.format_exc()) 
            return templates.TemplateResponse('default_error.html', context={
                'request': request,
                'status_code': 500,
                'detail': 'Internal server error'
            }, status_code=500)
