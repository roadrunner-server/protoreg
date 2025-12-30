package protoreg

import (
	er "errors"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/v2/protoresolve"
	"github.com/roadrunner-server/errors"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Registry interface {
	Registry() *protoresolve.Registry
	Services() map[string]*desc.ServiceDescriptor
	FindMethodByFullPath(method string) (*desc.MethodDescriptor, error)
}

type ProtoRegistry struct {
	registry *protoresolve.Registry
	services map[string]*desc.ServiceDescriptor
}

func (p *Plugin) InitRegistry() (*ProtoRegistry, error) {
	reg := &ProtoRegistry{
		services: make(map[string]*desc.ServiceDescriptor),
	}

	parser := &protoparse.Parser{
		ImportPaths: p.config.ImportPaths,
	}

	files := &protoregistry.Files{}

	fds, err := parser.ParseFiles(p.config.Proto...)
	if err != nil {
		return nil, err
	}

	for _, file := range fds {
		err = files.RegisterFile(file.UnwrapFile())
		if err != nil {
			return nil, err
		}

		err = registerDependencies(files, parser, file.GetDependencies())
		if err != nil {
			return nil, err
		}

		for _, service := range file.GetServices() {
			reg.services[service.GetFullyQualifiedName()] = service
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
func (reg *ProtoRegistry) Services() map[string]*desc.ServiceDescriptor {
	return reg.services
}

// FindMethodByFullPath finds a method descriptor by full method path
func (reg *ProtoRegistry) FindMethodByFullPath(method string) (*desc.MethodDescriptor, error) {
	parts := strings.Split(strings.TrimPrefix(method, "/"), "/")
	if len(parts) != 2 {
		return nil, errors.Errorf("Unexpected method")
	}

	service, ok := reg.services[parts[0]]
	if !ok {
		return nil, errors.Errorf("Service not found: %s", parts[0])
	}

	return service.FindMethodByName(parts[1]), nil
}

func registerDependencies(files *protoregistry.Files, parser *protoparse.Parser, deps []*desc.FileDescriptor) error {
	const op = errors.Op("protoreg_registry_parse_proto")

	if len(deps) == 0 {
		return nil
	}

	filenames := filenamesFromDesc(deps)

	fds, err := parser.ParseFiles(filenames...)
	if err != nil {
		return errors.E(op, err)
	}

	for _, fd := range fds {
		_, err = files.FindFileByPath(fd.GetName())
		if er.Is(err, protoregistry.NotFound) {
			err = files.RegisterFile(fd.UnwrapFile())
			if err != nil {
				return errors.E(op, err)
			}
		}

		err := registerDependencies(files, parser, fd.GetDependencies())
		if err != nil {
			return errors.E(op, err)
		}
	}

	return nil
}

func filenamesFromDesc(files []*desc.FileDescriptor) []string {
	names := make([]string, len(files))

	for i, file := range files {
		names[i] = file.GetName()
	}

	return names
}
