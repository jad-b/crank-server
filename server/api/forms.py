from django.forms import ModelForm

from .models import Workout, Exercise, Set


class WorkoutForm(ModelForm):
    class Meta:
        model = Workout
        fields = ['comment']


class ExerciseForm(ModelForm):
    class Meta:
        model = Exercise
        fields = ['name', 'comment', 'rating']


class SetForm(ModelForm):
    class Meta:
        model = Set
        fields = ['exercise', 'reps', 'weight', 'order', 'rest']



