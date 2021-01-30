# Copyright 2021 Northern.tech AS
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.

import uuid
import pytest

from common import management_api_with_params
from internal_api import InternalAPIClient


@pytest.fixture
def device_id():
    client = InternalAPIClient()
    device_id = str(uuid.uuid4())
    new_device = {"device_id": device_id}
    r = client.provision_device_with_http_info(
        tenant_id="", new_device=new_device, _preload_content=False
    )
    assert r.status == 201
    yield device_id
    r = client.decommission_device_with_http_info(
        tenant_id="", device_id=device_id, _preload_content=False
    )
    assert r.status == 204


class TestManagementConfig:
    def test_config_device(self, device_id):
        user_id = str(uuid.uuid4())
        client = management_api_with_params(user_id=user_id)
        #
        # get the configuration (empty)
        r = client.get_key_value_pair_store(device_id)
        data = r.to_dict()
        assert {
            "id": device_id,
            "actual": None,
            "expected": None,
            "report_ts": "0001-01-01T00:00:00Z",
        } == {k: data[k] for k in ("id", "actual", "expected", "report_ts")}
        #
        # set the initial configuration
        configuration = {
            "expected": [
                {"key": "key", "value": "value"},
                {"key": "another-key", "value": "another-value"},
                #  {"key": "empty-key", "value": ""},
                {"key": "dollar-key", "value": "$"},
            ],
        }
        r = client.store_a_set_of_key_value_pairs_with_http_info(
            device_id, inline_object=configuration, _preload_content=False
        )
        assert r.status == 201
        # data = r.to_dict()
        # assert {
        #     "id": device_id,
        #     "actual": None,
        #     "expected": None,
        #     "report_ts": "0001-01-01T00:00:00Z",
        # } == {k: data[k] for k in ("id", "actual", "expected", "report_ts")}
        #
        # get the configuration
        r = client.get_key_value_pair_store(device_id)
        data = r.to_dict()
        assert {
            "id": device_id,
            "actual": None,
            "expected": [
                {"key": "key", "value": "value"},
                {"key": "another-key", "value": "another-value"},
                {"key": "dollar-key", "value": "$"},
            ],
            "report_ts": "0001-01-01T00:00:00Z",
        } == {k: data[k] for k in ("id", "actual", "expected", "report_ts")}
        #
        # replace the configuration
        configuration = {
            "expected": [
                {"key": "key", "value": "update-value"},
                {"key": "additional-key", "value": '"'},
            ],
        }
        r = client.store_a_set_of_key_value_pairs_with_http_info(
            device_id, inline_object=configuration, _preload_content=False
        )
        assert r.status == 201
        #
        # get the configuration
        r = client.get_key_value_pair_store(device_id)
        data = r.to_dict()
        assert {
            "id": device_id,
            "actual": None,
            "expected": [
                {"key": "key", "value": "update-value"},
                {"key": "additional-key", "value": '"'},
            ],
            "report_ts": "0001-01-01T00:00:00Z",
        } == {k: data[k] for k in ("id", "actual", "expected", "report_ts")}
        #
        # remove the configuration
        configuration = {
            "expected": [],
        }
        r = client.store_a_set_of_key_value_pairs_with_http_info(
            device_id, inline_object=configuration, _preload_content=False
        )
        assert r.status == 201
        #
        # get the configuration
        r = client.get_key_value_pair_store(device_id)
        data = r.to_dict()
        assert {
            "id": device_id,
            "actual": None,
            "expected": [],
            "report_ts": "0001-01-01T00:00:00Z",
        } == {k: data[k] for k in ("id", "actual", "expected", "report_ts")}
