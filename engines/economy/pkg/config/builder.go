package config

import (
	"fmt"
	"westex/engines/economy/pkg/entities"
)

// BuildRegionFromConfig creates a Region from configuration
func BuildRegionFromConfig(config *RegionConfig) (*entities.Region, error) {
	region := entities.NewRegion(config.Region.Name)

	// Create problems map for lookup
	problemsMap := make(map[string]*entities.Problem)
	for _, pConfig := range config.Problems {
		problem := entities.NewProblem(pConfig.Name, pConfig.Description, pConfig.Demand)
		problem.IsBasicNeed = pConfig.IsBasicNeed
		region.AddProblem(problem)
		problemsMap[pConfig.Name] = problem
	}

	// Create resources map for lookup
	resourcesMap := make(map[string]*entities.Resource)
	for _, rConfig := range config.Resources {
		resource := entities.NewResource(rConfig.Name, rConfig.Unit)
		resource.Quantity = rConfig.InitialQuantity
		resource.IsFree = rConfig.IsFree
		resource.RegenerationRate = rConfig.RegenerationRate
		region.AddResource(resource)
		resourcesMap[rConfig.Name] = resource
	}

	// Create industries
	for _, iConfig := range config.Industries {
		// Get problems this industry solves
		solvedProblems := make([]*entities.Problem, 0)
		for _, problemName := range iConfig.SolvesProblems {
			if problem, exists := problemsMap[problemName]; exists {
				solvedProblems = append(solvedProblems, problem)
			} else {
				return nil, fmt.Errorf("industry %s references unknown problem: %s", iConfig.Name, problemName)
			}
		}

		// Get input resources
		inputResources := make([]*entities.Resource, 0)
		for _, resourceName := range iConfig.InputResources {
			if resource, exists := resourcesMap[resourceName]; exists {
				inputResources = append(inputResources, resource)
			} else {
				return nil, fmt.Errorf("industry %s references unknown input resource: %s", iConfig.Name, resourceName)
			}
		}

		// Create output resources (products)
		outputResources := make([]*entities.Resource, 0)
		for _, resourceName := range iConfig.OutputResources {
			// Check if resource already exists
			if resource, exists := resourcesMap[resourceName]; exists {
				outputResources = append(outputResources, resource)
			} else {
				// Create new product resource
				resource := entities.NewResource(resourceName, "units")
				resource.Quantity = 0 // Products start at 0
				outputResources = append(outputResources, resource)
				resourcesMap[resourceName] = resource
			}
		}

		// Create industry
		industry := entities.CreateIndustry(iConfig.Name).
			SetupIndustry(solvedProblems, inputResources, outputResources).
			UpdateLabor(iConfig.LaborNeeded).
			SetInitialCapital(iConfig.InitialCapital)

		region.AddIndustry(industry)
	}

	// Create population segments map
	segmentsMap := make(map[string]*entities.PopulationSegment)
	for _, sConfig := range config.Population.Segments {
		// Get problems for this segment
		segmentProblems := make([]*entities.Problem, 0)
		for _, problemName := range sConfig.HasProblems {
			if problem, exists := problemsMap[problemName]; exists {
				segmentProblems = append(segmentProblems, problem)
			}
		}

		size := int(float32(config.Population.TotalSize) * sConfig.Percentage)
		segment := &entities.PopulationSegment{
			Name:     sConfig.Name,
			Problems: segmentProblems,
			Size:     size,
		}
		segmentsMap[sConfig.Name] = segment
		region.AddPopulationSegment(segment)
	}

	// Create people
	personID := 1
	for _, sConfig := range config.Population.Segments {
		segment := segmentsMap[sConfig.Name]
		count := int(float32(config.Population.TotalSize) * sConfig.Percentage)

		for i := 0; i < count; i++ {
			person := entities.NewPerson(
				fmt.Sprintf("Person-%d", personID),
				sConfig.InitialMoney,
				sConfig.LaborHours,
			)
			person.AddSegment(segment)
			region.AddPerson(person)
			personID++
		}
	}

	return region, nil
}
