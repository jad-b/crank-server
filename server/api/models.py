from django.db import models


class Workout(models.Model):
    """Model representing an entire workout."""
    start = models.DateTimeField(unique=True)
    comment = models.CharField(max_length=256, blank=True)


class ExerciseName(models.Model):
    """Names of different Exercises."""
    name = models.CharField(max_length=64, unique=True)
    base = models.BooleanField(default=False,
                               help_text='Basic exercise, e.g. Squat, Press')
    modifiers = models.ManyToManyField('self',
                                       help_text='Allowed modifiers')


class Exercise(models.Model):
    workout = models.ForeignKey(Workout)
    name = models.ManyToManyField(ExerciseName,
                                  help_text='Name of the exercise')
    comment = models.CharField(max_length=256, blank=True)
    rating = models.PositiveSmallIntegerField(default=0,
                                              blank=True,
                                              help_text='Quality of exercise')


class Set(models.Model):
    """One set of an exercise."""
    exercise = models.ForeignKey(Exercise)
    reps = models.PositiveSmallIntegerField()
    weight = models.PositiveSmallIntegerField()
    order = models.PositiveSmallIntegerField()
    rest = models.PositiveSmallIntegerField(null=True, blank=True)
