import json

from django.contrib.auth.models import User
from django.test import SimpleTestCase
from django.core.urlresolvers import reverse


class TestLogin(SimpleTestCase):
    auth_url = reverse('api:login')

    @classmethod
    def setUpClass(self):
        User.objects.create_user('jdb', '', 'jdb')
        # Create an inactive user
        rumple = User.objects.create_user('rumplestiltskin', '', 'inactive')
        rumple.is_active = False
        rumple.save()

    def test_no_user(self):
        response = self.client.post(self.auth_url)

        self.assertEqual(response.status_code, 200)
        msg = json.loads(response.content.decode())
        self.assertTrue(msg['error'])
        self.assertEqual(msg['message'], 'Credentials are invalid')

    def test_bad_user(self):
        response = self.client.post(self.auth_url, {'username': 'bad',
                                                    'password': 'user'})

        self.assertEqual(response.status_code, 200)
        msg = json.loads(response.content.decode())
        self.assertTrue(msg['error'])
        self.assertEqual(msg['message'], 'Credentials are invalid')

    def test_inactive_user(self):
        response = self.client.post(
            self.auth_url,
            {'username': 'rumplestiltskin', 'password': 'inactive'})

        self.assertEqual(response.status_code, 200)
        msg = json.loads(response.content.decode())
        self.assertTrue(msg['error'])
        self.assertEqual(msg['message'], 'Your account is not active')

    def test_good_user(self):
        response = self.client.post(self.auth_url, {'username': 'jdb',
                                                    'password': 'jdb'})

        self.assertEqual(response.status_code, 200)
        msg = json.loads(response.content.decode())
        self.assertFalse(msg['error'])
        self.assertEqual(msg['message'], 'You have successfully authenticated')
