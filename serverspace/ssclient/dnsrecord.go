package ssclient

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	DNSRecordResponse struct {
		ID             int `json:"id,omitempty"`
		Name           string `json:"name,omitempty"`
		Type           string `json:"type,omitempty"`
		IP             string `json:"ip,omitempty"`
		MailHost       string `json:"mail_host,omitempty"`
		Priority       int    `json:"priority,omitempty"`
		CanonicalName  string `json:"canonical_name,omitempty"`
		NameServerHost string `json:"name_server_host,omitempty"`
		Text           string `json:"text,omitempty"`
		Protocol       string `json:"protocol,omitempty"`
		Service        string `json:"service,omitempty"`
		Weight         int    `json:"weight,omitempty"`
		Port           int    `json:"port,omitempty"`
		Target         string `json:"target,omitempty"`
		TTL            string `json:"ttl,omitempty"`
	}

	dnsRecordResponseWrap struct {
		Record *DNSRecordResponse `json:"record,omitempty"`
	}
)

func (c *SSClient) GetDNSRecord(domain string, record_id string) (*DNSRecordResponse, error) {
	url := fmt.Sprintf("%s/%s/records/%s", domainBaseURL, domain, record_id)
	resp, err := makeRequest(c.client, url, methodGet, nil, &dnsRecordResponseWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*dnsRecordResponseWrap).Record, nil
}

func (c *SSClient) CreateDNSRecord(
	name string,
	rtype string,
	ip string,
	mailHost string,
	priority int,
	canonicalName string,
	nameServerHost string,
	text string,
	protocol string,
	service string,
	weight int,
	port int,
	target string,
	ttl string) (*TaskIDWrap, error) {
	payload := map[string]interface{}{
		"name":             name,
		"type":             rtype,
		"ip":               ip,
		"mail_host":        mailHost,
		"priority":         priority,
		"canonical_name":   canonicalName,
		"name_server_host": nameServerHost,
		"text":             text,
		"protocol":         protocol,
		"service":          service,
		"weight":           weight,
		"port":             port,
		"target":           target,
		"ttl":              ttl,
	}
	domain := strings.Join(strings.Split(name, ".")[1:],".")
	url := fmt.Sprintf("%s/%s/records", domainBaseURL, domain)
	resp, err := makeRequest(c.client, url, methodPost, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) CreateDNSRecordAndWait(
	name string,
	rtype string,
	ip string,
	mailHost string,
	priority int,
	canonicalName string,
	nameServerHost string,
	text string,
	protocol string,
	service string,
	weight int,
	port int,
	target string,
	ttl string) (*DNSRecordResponse, error) {
	taskWrap, err := c.CreateDNSRecord(name, rtype, ip, mailHost, priority, canonicalName, nameServerHost, text,
		protocol, service, weight, port, target, ttl)
	if err != nil {
		return nil, err
	}
	return c.waitDNSRecord(taskWrap.ID)
}

func (c *SSClient) UpdateDNSRecord(recordID string,
	name string,
	rtype string,
	ip string,
	mailHost string,
	priority int,
	canonicalName string,
	nameServerHost string,
	text string,
	protocol string,
	service string,
	weight int,
	port int,
	target string,
	ttl string) (*TaskIDWrap, error) {
	domain := strings.Join(strings.Split(name, ".")[1:], ".")
	payload := map[string]interface{}{
		"name":             name,
		"type":             rtype,
		"ip":               ip,
		"mail_host":        mailHost,
		"priority":         priority,
		"canonical_name":   canonicalName,
		"name_server_host": nameServerHost,
		"text":             text,
		"protocol":         protocol,
		"service":          service,
		"weight":           weight,
		"port":             port,
		"target":           target,
		"ttl":              ttl,
	}
	url := fmt.Sprintf("%s/%s/records/%s", domainBaseURL, domain, recordID)
	resp, err := makeRequest(c.client, url, methodPut, payload, &TaskIDWrap{})
	if err != nil {
		return nil, err
	}
	return resp.(*TaskIDWrap), nil
}

func (c *SSClient) UpdateDNSRecordAndWait(recordID string,
	name string,
	rtype string,
	ip string,
	mailHost string,
	priority int,
	canonicalName string,
	nameServerHost string,
	text string,
	protocol string,
	service string,
	weight int,
	port int,
	target string,
	ttl string) (*DNSRecordResponse, error) {
	taskWrap, err := c.UpdateDNSRecord(recordID, name, rtype, ip, mailHost, priority, canonicalName, nameServerHost, text, protocol, service, weight, port, target, ttl)
	if err != nil {
		return nil, err
	}
	return c.waitDNSRecord(taskWrap.ID)
}

func (c *SSClient) DeleteDNSRecord(domain string, recordID string) error {
	url := fmt.Sprintf("%s/%s/records/%s", domainBaseURL, domain, recordID)
	_, err := makeRequest(c.client, url, methodDelete, nil, nil)
	return err
}

func (c *SSClient) waitDNSRecord(taskID string) (*DNSRecordResponse, error) {
	record, err := c.waitTaskCompletion(taskID)
	if err != nil {
		return nil, err
	}
	return c.GetDNSRecord(record.DomainID, strconv.Itoa(record.RecordID))
}
