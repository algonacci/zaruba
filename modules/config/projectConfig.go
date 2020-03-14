package config

import (
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// ProjectConfig configuration
type ProjectConfig struct {
	dirName      string
	projectName  string
	environments *Environments
	components   map[string]*Component
	executions   []string
	links        map[string][]string
}

// GetProjectName get name of project
func (p *ProjectConfig) GetProjectName() (projectName string) {
	return p.projectName
}

// GetEnvironments get environments of project
func (p *ProjectConfig) GetEnvironments() (environments *Environments) {
	return p.environments
}

// GetComponents get components of project
func (p *ProjectConfig) GetComponents() (components map[string]*Component) {
	return p.components
}

// GetComponentByName get component of project by name
func (p *ProjectConfig) GetComponentByName(name string) (component *Component) {
	return p.components[name]
}

// GetExecutions get executions order of projects
func (p *ProjectConfig) GetExecutions() (executions []string) {
	return p.executions
}

// GetLinks get links in the project
func (p *ProjectConfig) GetLinks() (links map[string][]string) {
	return p.links
}

// GetLinkDestinationList get link by source
func (p *ProjectConfig) GetLinkDestinationList(source string) (destinationList []string) {
	return p.links[source]
}

// ToYaml get yaml representation of projectConfig
func (p *ProjectConfig) ToYaml() (str string, err error) {
	pYaml := &ProjectConfigYaml{
		ProjectName: p.GetProjectName(),
		Environments: EnvironmentsYaml{
			General:  p.GetEnvironments().general,
			Services: p.GetEnvironments().services,
		},
		Components: map[string]ComponentYaml{},
		Executions: p.GetExecutions(),
		Links:      p.GetLinks(),
	}
	for componentName, component := range p.GetComponents() {
		pYaml.Components[componentName] = ComponentYaml{
			Type:          component.GetType(),
			Origin:        component.GetOrigin(),
			Branch:        component.GetBranch(),
			Location:      component.GetLocation(),
			Start:         component.GetStartCommand(),
			Run:           component.GetRunCommand(),
			ContainerName: component.GetContainerName(),
		}
	}

	d, err := yaml.Marshal(*pYaml)
	if err != nil {
		return str, err
	}
	str = string(d)
	return str, err
}

// GetSortedLinkSources get sorted link sources
func (p *ProjectConfig) GetSortedLinkSources() (sortedSources []string) {
	sortedSources = []string{}
	for source := range p.links {
		sortedSources = append(sortedSources, source)
	}
	// sort keys
	sort.SliceStable(sortedSources, func(i int, j int) bool {
		firstSource, secondSource := sortedSources[i], sortedSources[j]
		// get destination
		firstDestinations := p.links[firstSource]
		// compare
		for _, destination := range firstDestinations {
			if strings.HasPrefix(destination, secondSource) {
				return true
			}
		}
		return false
	})
	return sortedSources
}

// GetSubrepoPrefixMap get map of all eligible component's subrepoPrefix (for git subtree)
func (p *ProjectConfig) GetSubrepoPrefixMap(projectDir string) (subRepoPrefixMap map[string]string) {
	subRepoPrefixMap = map[string]string{}
	for componentName, component := range p.components {
		location := component.location
		origin := component.origin
		branch := component.branch
		if location == "" || origin == "" || branch == "" {
			continue
		}
		subRepoPrefix := getSubrepoPrefix(projectDir, location)
		subRepoPrefixMap[componentName] = subRepoPrefix
	}
	return subRepoPrefixMap
}

func getSubrepoPrefix(projectDir, location string) string {
	if !strings.HasPrefix(location, projectDir) {
		return location
	}
	return strings.Trim(strings.TrimPrefix(location, projectDir), string(os.PathSeparator))
}

func (p *ProjectConfig) fromProjectConfigYaml(pYaml *ProjectConfigYaml, directory string) *ProjectConfig {
	// load pYaml into p
	p.dirName = directory
	p.projectName = pYaml.ProjectName
	p.environments = &Environments{
		general:  pYaml.Environments.General,
		services: pYaml.Environments.Services,
	}
	p.components = make(map[string]*Component)
	p.executions = pYaml.Executions
	p.links = pYaml.Links
	for componentName, component := range pYaml.Components {
		p.components[componentName] = &Component{
			componentType: component.Type,
			origin:        component.Origin,
			branch:        component.Branch,
			location:      component.Location,
			start:         component.Start,
			run:           component.Run,
			containerName: component.ContainerName,
		}
	}
	return p
}

// Environments describe environment variables in general and for each services
type Environments struct {
	general  map[string]string
	services map[string]map[string]string
}

// GetGeneralVariables get general environment variables
func (e *Environments) GetGeneralVariables() (general map[string]string) {
	return e.general
}

// GetAllServicesVariables get all service variables (as map)
func (e *Environments) GetAllServicesVariables() (services map[string]map[string]string) {
	return e.services
}

// GetServiceVariables get variable of a service
func (e *Environments) GetServiceVariables(serviceName string) (variables map[string]string) {
	variables = map[string]string{}
	for key, val := range e.general {
		variables[key] = val
	}
	if serviceEnv, exists := e.services[serviceName]; exists {
		for key, val := range serviceEnv {
			variables[key] = val
		}
	}
	return variables
}

// Component describe component specs
type Component struct {
	componentType string
	origin        string
	branch        string
	location      string
	start         string
	run           string
	containerName string
}

// GetType get component type
func (c *Component) GetType() (componentType string) {
	return c.componentType
}

// GetOrigin get component origin
func (c *Component) GetOrigin() (origin string) {
	return c.origin
}

// GetBranch get component branch
func (c *Component) GetBranch() (branch string) {
	return c.branch
}

// GetLocation get component location
func (c *Component) GetLocation() (location string) {
	return c.location
}

// GetStartCommand get component start command
func (c *Component) GetStartCommand() (start string) {
	return c.start
}

// GetRunCommand get component run command
func (c *Component) GetRunCommand() (run string) {
	return c.run
}

// GetContainerName get component container name
func (c *Component) GetContainerName() (containerName string) {
	return c.containerName
}
