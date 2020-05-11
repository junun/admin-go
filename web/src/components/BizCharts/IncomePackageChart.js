import React, {PureComponent} from "react";
import {Chart, Geom, Tooltip, Coord, Label, Legend} from 'bizcharts';
import {Card, Spin} from "antd";
import DataSet from "@antv/data-set";


class IncomePackageChart extends PureComponent {
  componentDidMount() {
    const e = document.createEvent("Event");
    e.initEvent("resize", true, true);
    window.dispatchEvent(e);
  }

  render() {
    const { loading, title, data } = this.props;
    const { DataView } = DataSet;
    const dv = new DataView();
    const mapedData = data.map(x => {
      return {
        ...x,
        packageName: (x.packagePrice / 100).toFixed(2) + '元'
      };
    });
    dv.source(mapedData).transform({
      type: "percent",
      field: "quantity",
      dimension: "packagePrice",
      as: "percent"
    });
    const cols = {
      percent: {
        formatter: val => {
          val = (val * 100).toFixed(2) + "%";
          return val;
        }
      }
    };
    let chart;
    if (!data || data.length === 0) {
      chart = <Chart placeholder height={380} />;
    } else {
      chart = <Chart
        height={415}
        data={dv}
        scale={cols}
        padding={[10, 50, 30, 50]}
        forceFit
      >
        <Coord type="theta" radius={0.6} />
        <Legend
          position="bottom"
          offsetY={-30}
        />
        <Tooltip
          showTitle={false}
          itemTpl="<li><span style=&quot;background-color:{color};&quot; class=&quot;g2-tooltip-marker&quot;></span>{name}: {value}</li>"
        />
        <Geom
          type="intervalStack"
          position="quantity"
          color="packageName"
          tooltip={[
            "packagePrice*quantity",
            (packagePrice, quantity) => {
              // percent = (percent * 100).toFixed(2) + "%";
              return {
                name: (packagePrice / 100).toFixed(2) + '元',
                value: quantity + '盏',
              };
            }
          ]}
          style={{
            lineWidth: 1,
            stroke: "#fff"
          }}
        >
          <Label
            content="percent"
            htmlTemplate={(val, item) => {
              const label = (item.point.packagePrice / 100).toFixed(2);
              const labelHTML = '<span style="color: #999; font-size: 14px; display: inline-block; white-space: nowrap">' + label + '元</span>';
              const valueHTML = '<span style="color: #333; font-size: 16px;">' + val + '</span>';
              return labelHTML + '<br>' + valueHTML;
            }}
          />
        </Geom>
      </Chart>;
    }

    return <Card title={title} style={{height: '492px'}} bodyStyle={{ padding: '10px'}}>
      {loading ?
        <Spin><div style={{ height: '360px'}}></div></Spin>
        :
        chart
      }
    </Card>;
  }
}

export default IncomePackageChart;
