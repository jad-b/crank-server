#!/usr/bin/env python
import argparse
import logging
import yaml

LOGGER = logging.getLogger(__name__)


def extract_includes(filepath='main.yaml'):
    """Given the provided *filepath*, load the list of swaggregator include
    statements, in the order that they appear in *filepath*

    :param filepath: The path to the swaggregator compliant yaml file
    :return: A :const:`list` of swagger files to include
    """
    includes = []
    include_stmt = '#include:'
    with open(filepath) as f:
        comment_lines = (l for l in f if l.startswith('#'))
        for line in comment_lines:
            if line.startswith(include_stmt):
                includes.append(line.replace(include_stmt, '').strip())
    return includes


class SwaggerSpec:
    """A object representation of a swagger API specification namespace. Adds
    support to the swagger spec language to add the use of "include" statements
    in your swagger spec file.
    ```
    ---
    #include: other_package/swagger.yaml
    swagger: '2.0'
    ```
    Included specification files will have an "in-order" precedence (ie, the
    last include statement could override attributes of the same name included
    from spec files that were included earlier).

    Currently, there is support for merging `tags`, `schemes`, `produces`,
    `parameters`, `paths`, and `definitions` from included spec files into the
    namespace of a SwaggerSpec instance
    """

    # Assortment of attributes to attempt to import into this specifications
    # namespace when loading included spec files.
    __attrs__ = [
        ('tags', list),
        ('schemes', list),
        ('produces', list),
        ('parameters', dict),
        ('paths', dict),
        ('definitions', dict)
    ]

    def __init__(self, spec='main.yaml'):
        """Create a new SwaggerSpec instance based off of the root
        specification from the file specified by *spec*

        :param spec: The swagger spec file entry point
        """
        super(SwaggerSpec, self).__init__()
        self._spec = spec

        # Set our attributes to safe defaults in case the main swagger spec
        # doesn't specify them
        for key, typ in self.__attrs__:
            setattr(self, key, typ())

        raw_data = self.load(spec)
        # Load our main swagger file's specification into this instance
        for k, v, in raw_data.items():
            setattr(self, k, v)

    def load_includes(self):
        """Load include statements and register them into this namespace"""
        includes = extract_includes(self._spec)
        for include in includes:
            self._include(include)

    def _include(self, include):
        """Recursive include method used to recursively include any
        specification files
        """
        self.merge(include)
        includes = extract_includes(include)
        for include in includes:
            self._include(include)

    def merge(self, include):
        """Merge the data stored in the *include* spec file, into our current
        SwaggerSpec instance
        """
        data = self.load(include)
        for key, typ in self.__attrs__:
            if typ is dict:
                self.__dict__[key].update(data.get(key, typ()))
            else:
                self.__dict__[key] += data.get(key, typ())

    def load(self, filepath='main.yaml'):
        """Load the main swaggregator file entry point

        :param filepath: The path to the main swaggregator entry point
        :return: The loaded yaml contents of *filepath* as a :const:`dict`
        """
        with open(filepath) as f:
            return yaml.load(f)

    def output(self, filepath):
        """Write the contents of this SwaggerSpec instance to *filepath*

        :param filepath: The path to the file to output the fully compiled spec
            out to
        """
        with open(filepath, 'w') as f:
            f.write(self.__yaml__())

    def __yaml__(self):
        """Express this SwaggerSpec instance as a YAML string"""
        data = {k: self.__dict__[k] for k in self.__dict__
                if not k.startswith('_')}
        return yaml.dump(data, default_flow_style=False)


def parse_args():
    """Parse commandline arguments and return the namespace they are compiled
    in by the ArgumentParser
    """
    parser = argparse.ArgumentParser(
        description='Aggregate swagger spec files rom multiple, local packages'
    )

    parser.add_argument('-m', '--main', action='store', default='main.yaml',
                       help='The main swagger file to load.')
    parser.add_argument('-o', '--out', action='store', default='output.yaml',
                       help='The name of the yaml file to output to.')
    return parser.parse_args()


def main():
    """Main function. Parses commandline arguments, creates a SwaggerSpec
    instance, loads all of the included spec files, and then outputs
    """
    args = parse_args()
    spec = SwaggerSpec(args.main)
    spec.load_includes()
    spec.output(args.out)


if __name__ == '__main__':
    main()

