from .project import MainProject, ServiceProject
import os

def test_service_project_generate():
    dir_name = './playground/service_project'
    try:
        os.removedirs(dir_name)
    except OSError:
        pass
    main_project = MainProject()
    main_project.generate(dir_name)
    # generate service project
    service_project = ServiceProject()
    service_project.load_from_template('./test_resources/service.zaruba.yaml')
    service_project.generate(dir_name=dir_name, service_name='myService', image_name='myImage', container_name='myContainer', location='./test_resources/app', ports=[])
    # reload
    generated_project = ServiceProject()
    generated_project.load(dir_name=dir_name, service_name='myService')
    # assert generated project
    assert len(generated_project.get(['includes'])) == 1
    assert generated_project.get(['includes', 0]) == '${ZARUBA_HOME}/scripts/core.zaruba.yaml'
    # runMyService
    assert generated_project.get(['tasks', 'runMyService', 'extend']) == 'core.startService'
    assert generated_project.get(['tasks', 'runMyService', 'location']) == '.././test_resources/app'
    assert generated_project.get(['tasks', 'runMyService', 'configRef']) == 'myService'
    assert generated_project.get(['tasks', 'runMyService', 'envRef']) == 'myService'
    assert generated_project.get(['tasks', 'runMyService', 'lconfRef']) == 'myService'
    # buildMyService
    assert generated_project.get(['tasks', 'buildMyServiceImage', 'extend']) == 'core.buildDockerImage'
    assert generated_project.get(['tasks', 'buildMyServiceImage', 'configRef']) == 'myServiceContainer'
    # config
    assert generated_project.get(['configs', 'myServiceContainer', 'containerName']) == 'myContainer'
    assert generated_project.get(['configs', 'myServiceContainer', 'imageName']) == 'myimage'
    # lconfig
    assert generated_project.get(['lconfigs', 'myService', 'ports', 0]) == '{{ .GetEnv "PORT" }}'
    # envs
    assert generated_project.get(['envs', 'myService', 'PORT', 'from']) == 'MYSERVICE_PORT'
    assert generated_project.get(['envs', 'myService', 'PORT', 'default']) == '3000'
    # assert main project
    main_project = generated_project.main_project
    assert len(main_project.get(['includes'])) == 2
    assert main_project.get(['includes', 0]) == '${ZARUBA_HOME}/scripts/core.zaruba.yaml'
    assert main_project.get(['includes', 1]) == 'zaruba-tasks/myService.zaruba.yaml'
    assert main_project.get(['tasks', 'run', 'dependencies', 0]) == 'runMyService'
    assert main_project.get(['tasks', 'runContainer', 'dependencies', 0]) == 'runMyServiceContainer'
    assert main_project.get(['tasks', 'stopContainer', 'dependencies', 0]) == 'stopMyServiceContainer'
    assert main_project.get(['tasks', 'removeContainer', 'dependencies', 0]) == 'removeMyServiceContainer'
    assert main_project.get(['tasks', 'buildImage', 'dependencies', 0]) == 'buildMyServiceImage'
    assert main_project.get(['tasks', 'pushImage', 'dependencies', 0]) == 'pushMyServiceImage'