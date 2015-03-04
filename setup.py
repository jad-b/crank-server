import os
from setuptools import setup


with open(os.path.join(os.path.dirname(__file__), 'README.rst')) as readme:
    README = readme.read()


def read_version(filename='version'):
    with open(filename) as f:
        return f.read().strip()

setup(
    name='crank-server',
    version=read_version(),
    packages=['crank'],
    include_package_data=True,
    license='MIT License',
    description='API server for Crank interactions',
    long_description=README,
    url='http://www.github.com/jad-b/crank-server',
    author='Jeremy Dobbins-Bucklad',
    author_email='j.american.db@gmail.com',
    classifiers=[
        'Environment :: Web Environment',
        'Framework :: Django',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.4',
        'Topic :: Internet :: WWW/HTTP',
        'Topic :: Internet :: WWW/HTTP :: Dynamic Content',
    ],
)
