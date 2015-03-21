# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Exercise',
            fields=[
                ('id', models.AutoField(serialize=False, primary_key=True, verbose_name='ID', auto_created=True)),
                ('comment', models.CharField(blank=True, max_length=256)),
                ('rating', models.PositiveSmallIntegerField(help_text='Quality of exercise', default=0, blank=True)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='ExerciseName',
            fields=[
                ('id', models.AutoField(serialize=False, primary_key=True, verbose_name='ID', auto_created=True)),
                ('name', models.CharField(unique=True, max_length=64)),
                ('base', models.BooleanField(help_text='Basic exercise, e.g. Squat, Press', default=False)),
                ('modifiers', models.ManyToManyField(help_text='Allowed modifiers', related_name='modifiers_rel_+', to='api.ExerciseName')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Set',
            fields=[
                ('id', models.AutoField(serialize=False, primary_key=True, verbose_name='ID', auto_created=True)),
                ('reps', models.PositiveSmallIntegerField()),
                ('weight', models.PositiveSmallIntegerField()),
                ('order', models.PositiveSmallIntegerField()),
                ('rest', models.PositiveSmallIntegerField(blank=True, null=True)),
                ('exercise', models.ForeignKey(to='api.Exercise')),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.CreateModel(
            name='Workout',
            fields=[
                ('id', models.AutoField(serialize=False, primary_key=True, verbose_name='ID', auto_created=True)),
                ('start', models.DateTimeField(unique=True)),
                ('comment', models.CharField(blank=True, max_length=256)),
            ],
            options={
            },
            bases=(models.Model,),
        ),
        migrations.AddField(
            model_name='exercise',
            name='name',
            field=models.ManyToManyField(help_text='Name of the exercise', to='api.ExerciseName'),
            preserve_default=True,
        ),
        migrations.AddField(
            model_name='exercise',
            name='workout',
            field=models.ForeignKey(to='api.Workout'),
            preserve_default=True,
        ),
    ]
