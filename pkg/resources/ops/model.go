package ops

import (
	"fmt"
	"strings"

	"k8s.io/kops/pkg/resources"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup"
)

type deletionAdapter struct {
	Task fi.Task
}

// ModelAdapter adapts model-based deletions to resources.
func ModelAdapter(cloud fi.Cloud, builder ...fi.HasDeletions) (map[string]*resources.Resource, error) {
	l := &cloudup.Loader{}
	l.Init()
	l.Builders = make([]fi.ModelBuilder, 0, len(builder))
	for _, b := range builder {
		l.Builders = append(l.Builders, b)
	}
	deletions, err := l.FindDeletions(cloud, map[string]fi.Lifecycle{})
	if err != nil {
		return nil, err
	}

	// todo: put the following in groupdeleter? how much of it is necessary
	context, err := fi.NewContext(target, cluster, cloud, keyStore, secretStore, configBase, checkExisting, c.TaskMap)
	if err != nil {
		return fmt.Errorf("error building context: %v", err)
	}
	defer context.Close()

	err = context.RunTasks(options)
	if err != nil {
		return fmt.Errorf("error running tasks: %v", err)
	}

	err = context.RunTasks(options)
	if err != nil {
		return fmt.Errorf("error running tasks: %v", err)
	}

	resourceMap := make(map[string]*resources.Resource, len(deletions))
	for taskName, task := range deletions {
		split := strings.SplitN(taskName, "/", 2)
		// todo call task.Find() at this point? put both in Obj?
		resourceMap[taskName] = &resources.Resource{
			Name: split[1],
			Type: split[0],
			ID: split[1],
			Deleter: deletionAdapter,
			Dumper: ,
			Obj: task,
		}
	}
	return resourceMap, nil
}

func deletionAdapter(cloud fi.Cloud, tracker *resources.Resource) error {
	task := tracker.Obj.(fi.Task)
	task.Find
}