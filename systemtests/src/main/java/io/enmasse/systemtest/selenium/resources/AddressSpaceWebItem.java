/*
 * Copyright 2018, EnMasse authors.
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */
package io.enmasse.systemtest.selenium.resources;

import org.openqa.selenium.By;
import org.openqa.selenium.WebElement;

public class AddressSpaceWebItem extends WebItem {

    private WebElement checkBox;
    private String name;
    private String namespace;
    private WebElement consoleRoute;
    private String type;
    private String status;
    private String created;
    private WebElement actionDropDown;

    public AddressSpaceWebItem(WebElement item) {
        this.webItem = item;
        this.checkBox = item.findElement(By.className("pf-c-table__check")).findElement(By.tagName("input"));
        this.name = parseName(item.findElement(By.xpath("//td[@data-label='Name/Namespace']")));
        this.namespace = parseNamespace(item.findElement(By.xpath("//td[@data-label='Name/Namespace']")));
        this.consoleRoute = parseRoute(item.findElement(By.xpath("//td[@data-label='Name/Namespace']")));
        this.type = item.findElement(By.xpath("//td[@data-label='Type']")).getText();
        this.status = item.findElement(By.xpath("//td[@data-label='Status']")).getText();
        this.created = item.findElement(By.xpath("//td[@data-label='Time created']")).getText();
        this.actionDropDown = item.findElement(By.className("pf-c-dropdown"));
    }

    public WebElement getCheckBox() {
        return checkBox;
    }

    public String getName() {
        return name;
    }

    public String getNamespace() {
        return namespace;
    }

    public WebElement getConsoleRoute() {
        return consoleRoute;
    }

    public String getType() {
        return type;
    }

    public String getStatus() {
        return status;
    }

    public String getCreated() {
        return created;
    }

    public WebElement getActionDropDown() {
        return actionDropDown;
    }

    private String parseName(WebElement elem) {
        try {
            return elem.findElement(By.tagName("a")).getText();
        } catch (Exception ex) {
            return elem.findElements(By.tagName("p")).get(0).getText();
        }
    }

    private WebElement parseRoute(WebElement elem) {
        try {
            return elem.findElement(By.tagName("a"));
        } catch (Exception ex) {
            return null;
        }
    }

    private String parseNamespace(WebElement elem) {
        if (elem.findElements(By.tagName("p")).size() > 1) {
            return elem.findElements(By.tagName("p")).get(1).getText();
        }
        return elem.findElement(By.tagName("p")).getText();
    }

    @Override
    public String toString() {
        return String.format("name: %s, namespace: %s, type: %s, status: %s, time created: %s",
                this.name,
                this.namespace,
                this.type,
                this.status,
                this.created);
    }
}
