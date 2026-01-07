package protoreg

import (
	"context"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/jhump/protoreflect/v2/protoresolve"
	"github.com/roadrunner-server/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Registry interface {
	Registry() *protoresolve.Registry
	Services() map[string]protoreflect.ServiceDescriptor
	FindMethodByFullPath(method string) (protoreflect.MethodDescriptor, error)
}

type ProtoRegistry struct {
	registry *protoresolve.Registry
	services map[string]protoreflect.ServiceDescriptor
}

func (p *Plugin) InitRegistry() (*ProtoRegistry, error) {
	reg := &ProtoRegistry{
		services: make(map[string]protoreflect.ServiceDescriptor),
	}

	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(&protocompile.SourceResolver{
			ImportPaths: p.config.ProtoPath,
		}),
	}

	ctx := context.Background()
	fds, err := compiler.Compile(ctx, p.config.Files...)
	if err != nil {
		return nil, err
	}

	files := &protoregistry.Files{}
	for _, fd := range fds {
		err = files.RegisterFile(fd)
		if err != nil {
			return nil, err
		}

		// Collect services
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			reg.services[string(svc.FullName())] = svc
		}
	}

	reg.registry, err = protoresolve.FromFiles(files)
	if err != nil {
		return nil, err
	}

	return reg, nil
}

// Registry returns the underlying registry
func (reg *ProtoRegistry) Registry() *protoresolve.Registry {
	return reg.registry
}

// Services returns the service descriptors map
func (reg *ProtoRegistry) Services() map[string]protoreflect.ServiceDescriptor {
	return reg.services
}

// FindMethodByFullPath finds a method descriptor by full method path
func (reg *ProtoRegistry) FindMethodByFullPath(method string) (protoreflect.MethodDescriptor, error) {
	parts := strings.Split(strings.TrimPrefix(method, "/"), "/")
	if len(parts) != 2 {
		return nil, errors.Errorf("Unexpected method")
	}

	service, ok := reg.services[parts[0]]
	if !ok {
		return nil, errors.Errorf("Service not found: %s", parts[0])
	}

	return service.Methods().ByName(protoreflect.Name(parts[1])), nil
}
