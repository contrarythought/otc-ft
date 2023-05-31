from selenium import webdriver
import time
import sys

def main(url, user_agent, file_path):
    # start web browser
    options = webdriver.ChromeOptions()
    options.add_argument("--headless=new")
    options.add_argument(f'--user-agent="{user_agent}"')
    browser=webdriver.Chrome(chrome_options=options)

    # get source code
    browser.get(url)
    html = browser.page_source
    time.sleep(2)

    # close web browser
    browser.close()
    browser.quit()

    file = open(file_path, "w")
    file.write(html)
    print(f"writing to {file_path}")
    file.close()

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("required: url, user-agent, file_path")
    else:
        main(sys.argv[1], sys.argv[2], sys.argv[3])