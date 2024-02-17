"""
Compile .yaml documents into .md files.
"""

import yaml
import json

response_meanings = {
    "200": "OK",
    "201": "Created",
    "204": "No Content",
    "400": "Bad Request",
    "401": "Unauthorized",
    "403": "Forbidden",
    "404": "Not Found",
    "500": "Internal Server Error",
}


def compile_yaml_to_md(yaml_path, md_path):
    with open(yaml_path, "r") as f:
        data = yaml.safe_load(f)

    md_content = ""

    def helper(level, d):
        nonlocal md_content

        for item in d:
            md_content += f"{'#' * level} {item}\n\n"
            if isinstance(d[item], dict):
                helper(level + 1, d[item])
            if isinstance(d[item], list):
                for x in d[item]:
                    method, path = x["api"].strip().split(" ", 1)
                    response = ""

                    for resp in x["response"]:
                        rescode, resjson = resp.strip().split(" ", 1)
                        resmean = response_meanings[rescode]
                        resjson = json.dumps(json.loads(resjson), indent=4)
                        response += f"```json\n{rescode} {resmean}\n{resjson}\n```\n\n"

                    params = "\n".join(
                        [
                            "|".join(
                                [
                                    y["name"],
                                    f' `{y["type"]}` ',
                                    y["desc"],
                                    f' `{y["required"]}` ',
                                ]
                            )
                            for y in x.get("param", [])
                        ]
                    )

                    markdown = f"""
{'#'*level} **{method}** _{path}_
> {x["desc"]}
{'#'*(level+1)} NEED AUTHENTICATION
{x["needAuth"]}
{'#'*(level+1)} PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
{params}
{'#'*(level+1)} EXAMPLE RESPONSE
{response}
"""

                    md_content += markdown

    helper(1, data)

    with open(md_path, "w") as f:
        f.write(md_content)


if __name__ == "__main__":
    compile_yaml_to_md("docs/apidocs-v1.yaml", "docs/apidocs-v1.md")