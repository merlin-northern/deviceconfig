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
        assert {"actual": None, "expected": None,} == {
            k: data[k] for k in ("actual", "expected")
        }
        #
        # set the initial configuration
        configuration = {
            "key": "value",
            "another-key": "another-value",
            #   "empty-key": "",
            "dollar-key": "$",
        }
        r = client.store_a_set_of_key_value_pairs_with_http_info(
            device_id, inline_object=configuration, _preload_content=False
        )
        assert r.status == 201
        # data = r.to_dict()
        # assert {
        #     "actual": None,
        #     "expected": None,
        # } == {k: data[k] for k in ("actual", "expected")}
        #
        # get the configuration
        r = client.get_key_value_pair_store(device_id)
        data = r.to_dict()
        assert {
            "actual": None,
            "expected": {
                "key": "value",
                "another-key": "another-value",
                "dollar-key": "$",
            },
        } == {k: data[k] for k in ("actual", "expected")}
        #
        # replace the configuration
        configuration = {
            "key": "update-value",
            "additional-key": '"',
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
            "actual": None,
            "expected": {"key": "update-value", "additional-key": '"',},
        } == {k: data[k] for k in ("actual", "expected")}
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
        assert {"actual": None, "expected": {},} == {
            k: data[k] for k in ("actual", "expected")
        }
