package scaleway

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayInstanceIPReverseDNS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayInstanceIPReverseDNSCreate,
		ReadContext:   resourceScalewayInstanceIPReverseDNSRead,
		UpdateContext: resourceScalewayInstanceIPReverseDNSUpdate,
		DeleteContext: resourceScalewayInstanceIPReverseDNSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultInstanceIPTimeout),
			Create:  schema.DefaultTimeout(defaultInstanceIPReverseDNSTimeout),
			Update:  schema.DefaultTimeout(defaultInstanceIPReverseDNSTimeout),
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			"ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP ID or IP address",
			},
			"reverse": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The reverse DNS for this IP",
			},
			"zone": zoneSchema(),
		},
		CustomizeDiff: customizeDiffLocalityCheck("ip_id"),
	}
}

func resourceScalewayInstanceIPReverseDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, err := instanceAPIWithZone(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetIP(&instance.GetIPRequest{
		IP:   expandID(d.Get("ip_id")),
		Zone: zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newZonedIDString(zone, res.IP.ID))

	if _, ok := d.GetOk("reverse"); ok {
		tflog.Debug(ctx, fmt.Sprintf("updating IP %q reverse to %q\n", d.Id(), d.Get("reverse")))

		updateReverseReq := &instance.UpdateIPRequest{
			Zone: zone,
			IP:   res.IP.ID,
		}

		if reverse, ok := d.GetOk("reverse"); ok {
			if isInstanceIPReverseResolved(ctx, instanceAPI, reverse.(string), d.Timeout(schema.TimeoutCreate), res.IP.ID, zone) {
				updateReverseReq.Reverse = &instance.NullableStringValue{Value: reverse.(string)}
			} else {
				return diag.FromErr(fmt.Errorf("your reverse must resolve. Ensure the command 'dig +short %s' matches your IP address ", reverse.(string)))
			}
		} else {
			updateReverseReq.Reverse = &instance.NullableStringValue{Null: true}
		}
		_, err = instanceAPI.UpdateIP(updateReverseReq, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayInstanceIPReverseDNSRead(ctx, d, meta)
}

func resourceScalewayInstanceIPReverseDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := instanceAPI.GetIP(&instance.GetIPRequest{
		IP:   ID,
		Zone: zone,
	}, scw.WithContext(ctx))
	if err != nil {
		// We check for 403 because instance API returns 403 for a deleted IP
		if is404Error(err) || is403Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("zone", string(zone))
	_ = d.Set("reverse", res.IP.Reverse)
	return nil
}

func resourceScalewayInstanceIPReverseDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("reverse") {
		tflog.Debug(ctx, fmt.Sprintf("updating IP %q reverse to %q\n", d.Id(), d.Get("reverse")))

		updateReverseReq := &instance.UpdateIPRequest{
			Zone: zone,
			IP:   ID,
		}

		if reverse, ok := d.GetOk("reverse"); ok {
			if isInstanceIPReverseResolved(ctx, instanceAPI, reverse.(string), d.Timeout(schema.TimeoutUpdate), ID, zone) {
				updateReverseReq.Reverse = &instance.NullableStringValue{Value: reverse.(string)}
			} else {
				return diag.FromErr(fmt.Errorf("your reverse must resolve. Ensure the command 'dig +short %s' matches your IP address ", reverse.(string)))
			}
		} else {
			updateReverseReq.Reverse = &instance.NullableStringValue{Null: true}
		}
		_, err = instanceAPI.UpdateIP(updateReverseReq, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayInstanceIPReverseDNSRead(ctx, d, meta)
}

func resourceScalewayInstanceIPReverseDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceAPI, zone, ID, err := instanceAPIWithZoneAndID(meta, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Unset the reverse dns on the IP
	updateReverseReq := &instance.UpdateIPRequest{
		Zone:    zone,
		IP:      ID,
		Reverse: &instance.NullableStringValue{Null: true},
	}
	_, err = instanceAPI.UpdateIP(updateReverseReq, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
