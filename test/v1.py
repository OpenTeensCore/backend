import requests
import yaml

docs = {}
APIRESULT = {}


def test_user_account(prefix: list):
    prefix.append("user")
    prefix.append("account")

    # Load API docs
    apis = docs
    for p in prefix[1:]:
        apis = apis[p]

    # Prepare session
    sess = requests.Session()

    # Test each API
    for api in apis:
        method, path = api["api"].split(" ", 1)
        url = prefix[0] + "/" + path
        print(f"[*] Testing [{method} {url}] ... ")

        # Prepare request
        data = {}
        if "param" in api:
            for param in api["param"]:
                if isinstance(param["example"], dict):
                    if "script" in param["example"]:
                        param["example"] = eval(param["example"]["script"])
                data[param["name"]] = param["example"]


        if api["needAuth"] is True:
            print("[*] \033[93mNeed Auth, skip\033[0m")
        else:
            # Send request
            if method == "GET":
                resp = sess.get(url, params=data)
            if method == "POST":
                resp = sess.post(url, data=data)
            if method == "PUT":
                resp = sess.put(url, data=data)
            if method == "DELETE":
                resp = sess.delete(url)

            # Recv
            res = resp.json()
            APIRESULT[api["api"]] = res
            
            # check response
            if resp.status_code == 200:
                print("[*] \033[92mOK\033[0m")
            else:
                print("[*] \033[91mFailed\033[0m")
                exit(-1)


def test_v1(prefix: list):
    prefix.append("v1")
    test_user_account(prefix)


if __name__ == "__main__":
    with open("docs/apidocs-v1.yaml", "r") as f:
        docs = yaml.load(f, Loader=yaml.FullLoader)

    test_v1(["http://localhost:5088"])
