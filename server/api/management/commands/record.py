import cmd

from django.core.management.base import BaseCommand


class Command(BaseCommand):
    """Command providing a CLI for entering workouts."""
    help = "Import Mailpipes model data from JSON into the database."

    def add_arguments(self, parser):
        return

    def handle(self, *args, **options):
        """Handle processing the commandline arguments and handling their
        individual cases
        """
        shell = EntryShell()
        try:
            shell.cmdloop()
        except KeyboardInterrupt:   # Accept C-D as exit
            shell.do_save(None)


class EntryShell(cmd.Cmd):
    """Command line shell for entering workouts by hand."""

    def __init__(self):
        super().__init__()

    def emptyline(self):
        self.do_help('')

    def do_EOF(self, args):
        raise KeyboardInterrupt

    def do_save(self, args):
        pass
